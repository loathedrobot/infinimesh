/*
Copyright © 2021-2023 Infinite Devices GmbH, Nikita Ivanovski info@slnt-opp.xyz

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package graph

import (
	"context"

	"github.com/arangodb/go-driver"
	inf "github.com/infinimesh/infinimesh/pkg/shared"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/infinimesh/infinimesh/pkg/graph/schema"
	"github.com/infinimesh/proto/node/access"
	"github.com/infinimesh/proto/node/namespaces"
	pb "github.com/infinimesh/proto/plugins"
)

type Plugin struct {
	*pb.Plugin
	driver.DocumentMeta
}

func (o *Plugin) ID() driver.DocumentID {
	return o.DocumentMeta.ID
}

func (o *Plugin) SetAccessLevel(level access.Level) {
	if o.Access == nil {
		o.Access = &access.Access{
			Level: level,
		}
		return
	}
	o.Access.Level = level
}

func NewBlankPluginDocument(key string) *Plugin {
	return &Plugin{
		Plugin: &pb.Plugin{
			Uuid: key,
		},
		DocumentMeta: NewBlankDocument(schema.PLUGINS_COL, key),
	}
}

type PluginsController struct {
	pb.UnimplementedPluginsServiceServer
	log *zap.Logger

	col     driver.Collection // Plugins Collection
	ns_ctrl *NamespacesController

	db driver.Database
}

func NewPluginsController(log *zap.Logger, db driver.Database) *PluginsController {
	ctx := context.TODO()
	col, _ := db.Collection(ctx, schema.PLUGINS_COL)
	return &PluginsController{
		log: log.Named("PluginsController"), col: col, db: db,
		ns_ctrl: NewNamespacesController(log, db),
	}
}

func ValidateRoot(ctx context.Context) bool {
	rootV := ctx.Value(inf.InfinimeshRootCtxKey)
	if rootV == nil {
		return false
	}

	root, ok := rootV.(bool)
	return ok && root
}

func ValidatePluginDocument(p *pb.Plugin) string {
	if p.Title == "" {
		return "Title cannot be empty"
	}

	if p.Kind == pb.PluginKind_UNKNOWN {
		return "Kind can't be Unknown"
	} else if p.Kind == pb.PluginKind_EMBEDDED && p.EmbeddedConf == nil {
		return "Kind is set to Embedded, but no conf provided"
	}

	return ""
}

func (c *PluginsController) Create(ctx context.Context, plug *pb.Plugin) (*pb.Plugin, error) {
	log := c.log.Named("Create")

	if !ValidateRoot(ctx) {
		return nil, status.Error(codes.PermissionDenied, "Not enough access rights to create Plugin")
	}

	log.Debug("Creating", zap.Any("plugin", plug))
	msg := ValidatePluginDocument(plug)
	if msg != "" {
		return nil, status.Error(codes.InvalidArgument, msg)
	}

	plugin := Plugin{Plugin: plug}
	meta, err := c.col.CreateDocument(ctx, plugin)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error while creating Plugin in DB: %v", err)
	}
	plugin.Uuid = meta.ID.Key()
	plugin.DocumentMeta = meta

	log.Debug("Created", zap.String("plugin", plugin.Uuid))
	return plugin.Plugin, nil
}

const listAllPluginsQuery = `FOR plug IN @@plugins RETURN plug`
const listAllPublicPluginsQuery = `
FOR plug IN @@plugins
FILTER plug.public
    RETURN plug
`

func (c *PluginsController) Get(ctx context.Context, plug *pb.Plugin) (*pb.Plugin, error) {
	log := c.log.Named("Get")

	meta, err := c.col.ReadDocument(ctx, plug.Uuid, plug)
	if err != nil {
		log.Warn("Couldn't get plugin", zap.Error(err))
		return nil, status.Error(codes.NotFound, "Plugin not found")
	}
	plug.Uuid = meta.ID.Key()

	if plug.Public || ValidateRoot(ctx) {
		return plug, nil
	}

	if plug.Namespace == nil || *plug.Namespace == "" {
		return nil, status.Error(codes.InvalidArgument, "Namespace is not given")
	}

	ns, err := c.ns_ctrl.Get(ctx, &namespaces.Namespace{Uuid: *plug.Namespace})
	if err != nil {
		return nil, err
	}

	if ns.Access.Level < access.Level_READ {
		return nil, status.Error(codes.PermissionDenied, "Not enough Access")
	}

	return plug, nil
}

func (c *PluginsController) List(ctx context.Context, r *pb.ListRequest) (*pb.Plugins, error) {
	log := c.log.Named("List")

	var cr driver.Cursor
	var err error

	if ValidateRoot(ctx) {
		cr, err = c.db.Query(ctx, listAllPluginsQuery, map[string]interface{}{
			"@plugins": schema.PLUGINS_COL,
		})

	} else if r.Namespace != nil && *r.Namespace != "" {
		cr, err = ListQuery(WithDepth(ctx, 1), log, c.db, NewBlankNamespaceDocument(*r.Namespace), schema.PLUGINS_COL)
	} else {
		cr, err = c.db.Query(ctx, listAllPublicPluginsQuery, map[string]interface{}{
			"@plugins": schema.PLUGINS_COL,
		})
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error getting Plugins from DB: %v", err)
	}
	defer cr.Close()

	var res []*pb.Plugin
	for {
		var plug pb.Plugin
		meta, err := cr.ReadDocument(ctx, &plug)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}
		plug.Uuid = meta.ID.Key()

		log.Debug("Got document", zap.Any("plugin", &plug))
		res = append(res, &plug)
	}

	return &pb.Plugins{
		Pool: res,
	}, nil
}

func (c *PluginsController) Update(ctx context.Context, plug *pb.Plugin) (*pb.Plugin, error) {
	log := c.log.Named("Update")

	if !ValidateRoot(ctx) {
		return nil, status.Error(codes.PermissionDenied, "Not enough access rights to update Plugin")
	}

	msg := ValidatePluginDocument(plug)
	if msg != "" {
		return nil, status.Error(codes.InvalidArgument, msg)
	}

	if plug.Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "No Plugin UUID has been provided")
	}

	plugin := Plugin{Plugin: plug}
	_, err := c.col.ReplaceDocument(ctx, plug.Uuid, plugin)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error while updating Plugin in DB: %v", err)
	}

	log.Debug("Updated", zap.Any("plugin", plugin))
	return plugin.Plugin, nil
}

func (c *PluginsController) Delete(ctx context.Context, plug *pb.Plugin) (*pb.Plugin, error) {
	if !ValidateRoot(ctx) {
		return nil, status.Error(codes.PermissionDenied, "Not enough access rights to delete Plugin")
	}

	_, err := c.col.RemoveDocument(ctx, plug.Uuid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error while deleting Plugin: %v", err)
	}

	return plug, nil
}
