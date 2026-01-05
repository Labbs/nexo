package document

import (
	fiberoapi "github.com/labbs/fiber-oapi"
	"github.com/labbs/nexo/application/document"
	"github.com/labbs/nexo/application/space"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/rs/zerolog"
)

type Controller struct {
	Config      config.Config
	Logger      zerolog.Logger
	FiberOapi   *fiberoapi.OApiGroup
	SpaceApp    *space.SpaceApp
	DocumentApp *document.DocumentApp
}

func SetupDocumentRouter(controller Controller) {
	// Search - must be before parameterized routes
	fiberoapi.Get(controller.FiberOapi, "/search", controller.SearchDocuments, fiberoapi.OpenAPIOptions{
		Summary:     "Search documents",
		Description: "Search documents by name or content",
		OperationID: "document.search",
		Tags:        []string{"Document", "Search"},
	})

	// Public document - before space routes
	fiberoapi.Get(controller.FiberOapi, "/public/:space_id/:identifier", controller.GetPublicDocument, fiberoapi.OpenAPIOptions{
		Summary:     "Get public document",
		Description: "Get a public document without authentication",
		OperationID: "document.getPublicDocument",
		Tags:        []string{"Document", "Public"},
	})

	// Space-level routes (no document ID)
	fiberoapi.Get(controller.FiberOapi, "/space/:space_id", controller.GetDocumentsFromSpace, fiberoapi.OpenAPIOptions{
		Summary:     "Get documents from space",
		Description: "Retrieve documents from a specific space",
		OperationID: "document.getDocumentsFromSpace",
		Tags:        []string{"Document"},
	})
	fiberoapi.Post(controller.FiberOapi, "/space/:space_id", controller.CreateDocument, fiberoapi.OpenAPIOptions{
		Summary:     "Create a new document",
		Description: "Create a new document in a specified space",
		OperationID: "document.createDocument",
		Tags:        []string{"Document"},
	})
	fiberoapi.Get(controller.FiberOapi, "/space/:space_id/trash", controller.GetTrash, fiberoapi.OpenAPIOptions{
		Summary:     "Get trash",
		Description: "Get all deleted documents in a space",
		OperationID: "document.getTrash",
		Tags:        []string{"Document", "Trash"},
	})

	// Routes with specific suffixes - MUST be before generic /:identifier routes
	// Version history
	fiberoapi.Get(controller.FiberOapi, "/space/:space_id/:document_id/versions", controller.ListVersions, fiberoapi.OpenAPIOptions{
		Summary:     "List document versions",
		Description: "Get version history for a document",
		OperationID: "document.listVersions",
		Tags:        []string{"Document", "Versions"},
	})
	fiberoapi.Get(controller.FiberOapi, "/space/:space_id/:document_id/versions/:version_id", controller.GetVersion, fiberoapi.OpenAPIOptions{
		Summary:     "Get document version",
		Description: "Get a specific version of a document",
		OperationID: "document.getVersion",
		Tags:        []string{"Document", "Versions"},
	})
	fiberoapi.Post(controller.FiberOapi, "/space/:space_id/:document_id/versions", controller.CreateVersion, fiberoapi.OpenAPIOptions{
		Summary:     "Create version snapshot",
		Description: "Manually create a version snapshot of the current document state",
		OperationID: "document.createVersion",
		Tags:        []string{"Document", "Versions"},
	})
	fiberoapi.Post(controller.FiberOapi, "/space/:space_id/:document_id/versions/:version_id/restore", controller.RestoreVersion, fiberoapi.OpenAPIOptions{
		Summary:     "Restore document version",
		Description: "Restore a document to a previous version",
		OperationID: "document.restoreVersion",
		Tags:        []string{"Document", "Versions"},
	})

	// Comments
	fiberoapi.Get(controller.FiberOapi, "/space/:space_id/:document_id/comments", controller.GetComments, fiberoapi.OpenAPIOptions{
		Summary:     "Get document comments",
		Description: "Get all comments for a document",
		OperationID: "document.getComments",
		Tags:        []string{"Document", "Comments"},
	})
	fiberoapi.Post(controller.FiberOapi, "/space/:space_id/:document_id/comments", controller.CreateComment, fiberoapi.OpenAPIOptions{
		Summary:     "Create comment",
		Description: "Create a new comment on a document",
		OperationID: "document.createComment",
		Tags:        []string{"Document", "Comments"},
	})
	fiberoapi.Put(controller.FiberOapi, "/space/:space_id/:document_id/comments/:comment_id", controller.UpdateComment, fiberoapi.OpenAPIOptions{
		Summary:     "Update comment",
		Description: "Update an existing comment",
		OperationID: "document.updateComment",
		Tags:        []string{"Document", "Comments"},
	})
	fiberoapi.Delete(controller.FiberOapi, "/space/:space_id/:document_id/comments/:comment_id", controller.DeleteComment, fiberoapi.OpenAPIOptions{
		Summary:     "Delete comment",
		Description: "Delete a comment",
		OperationID: "document.deleteComment",
		Tags:        []string{"Document", "Comments"},
	})
	fiberoapi.Patch(controller.FiberOapi, "/space/:space_id/:document_id/comments/:comment_id/resolve", controller.ResolveComment, fiberoapi.OpenAPIOptions{
		Summary:     "Resolve/unresolve comment",
		Description: "Mark a comment as resolved or unresolved",
		OperationID: "document.resolveComment",
		Tags:        []string{"Document", "Comments"},
	})

	// Document permissions
	fiberoapi.Get(controller.FiberOapi, "/space/:space_id/:document_id/permissions", controller.ListDocumentPermissions, fiberoapi.OpenAPIOptions{
		Summary:     "List document permissions",
		Description: "List all permissions for a specific document",
		OperationID: "document.listPermissions",
		Tags:        []string{"Document", "Permissions"},
	})
	fiberoapi.Put(controller.FiberOapi, "/space/:space_id/:document_id/permissions", controller.UpsertDocumentUserPermission, fiberoapi.OpenAPIOptions{
		Summary:     "Upsert document user permission",
		Description: "Create or update a user permission for a document",
		OperationID: "document.upsertUserPermission",
		Tags:        []string{"Document", "Permissions"},
	})
	fiberoapi.Delete(controller.FiberOapi, "/space/:space_id/:document_id/permissions/:user_id", controller.DeleteDocumentUserPermission, fiberoapi.OpenAPIOptions{
		Summary:     "Delete document user permission",
		Description: "Remove a user permission from a document",
		OperationID: "document.deleteUserPermission",
		Tags:        []string{"Document", "Permissions"},
	})

	// Other document-specific routes with suffixes
	fiberoapi.Post(controller.FiberOapi, "/space/:space_id/:document_id/restore", controller.RestoreDocument, fiberoapi.OpenAPIOptions{
		Summary:     "Restore document",
		Description: "Restore a deleted document from trash",
		OperationID: "document.restoreDocument",
		Tags:        []string{"Document", "Trash"},
	})
	fiberoapi.Put(controller.FiberOapi, "/space/:space_id/:document_id/public", controller.SetPublic, fiberoapi.OpenAPIOptions{
		Summary:     "Set document public status",
		Description: "Make a document public or private",
		OperationID: "document.setPublic",
		Tags:        []string{"Document", "Public"},
	})
	fiberoapi.Patch(controller.FiberOapi, "/space/:space_id/:id/move", controller.MoveDocument, fiberoapi.OpenAPIOptions{
		Summary:     "Move document",
		Description: "Move a document to a new parent (or root)",
		OperationID: "document.moveDocument",
		Tags:        []string{"Document"},
	})

	// Generic document routes - MUST be LAST
	fiberoapi.Get(controller.FiberOapi, "/space/:space_id/:identifier", controller.GetDocument, fiberoapi.OpenAPIOptions{
		Summary:     "Get document by ID",
		Description: "Retrieve a specific document by its ID",
		OperationID: "document.getDocument",
		Tags:        []string{"Document"},
	})
	fiberoapi.Put(controller.FiberOapi, "/space/:space_id/:id", controller.UpdateDocument, fiberoapi.OpenAPIOptions{
		Summary:     "Update document",
		Description: "Update a specific document",
		OperationID: "document.updateDocument",
		Tags:        []string{"Document"},
	})
	fiberoapi.Delete(controller.FiberOapi, "/space/:space_id/:identifier", controller.DeleteDocument, fiberoapi.OpenAPIOptions{
		Summary:     "Delete document",
		Description: "Delete a specific document",
		OperationID: "document.deleteDocument",
		Tags:        []string{"Document"},
	})
}
