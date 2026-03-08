package main

import (
	"net/http"
)

// openAPISpec is the OpenAPI 3.0 specification for the tickets API.
const openAPISpec = `
openapi: 3.0.3
info:
  title: Tickets API
  description: HTTP server for managing QR code tickets for a livestream event.
  version: 1.0.0
servers:
  - url: http://localhost:9000
paths:
  /qrcodes:
    post:
      summary: Create a ticket
      description: Creates a new QR code ticket. The server auto-generates the ID.
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/QRCodeInput'
      responses:
        '201':
          description: Ticket created
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/QRCode'
        '400':
          description: Invalid request body
        '500':
          description: Internal server error
    get:
      summary: List all tickets
      responses:
        '200':
          description: Array of QR code tickets
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/QRCode'
        '500':
          description: Internal server error
  /qrcodes/{id}:
    get:
      summary: Get a ticket by ID
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '200':
          description: Ticket found
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/QRCode'
        '404':
          description: Ticket not found
        '500':
          description: Internal server error
  /qrcodes/phone/{phone}:
    get:
      summary: Get a ticket by client phone number
      parameters:
        - name: phone
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Ticket found
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/QRCode'
        '404':
          description: Ticket not found
        '500':
          description: Internal server error
  /qrcodes/{id}/use:
    patch:
      summary: Mark a ticket as used
      description: Sets the used flag to true on the given ticket.
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '200':
          description: Updated ticket
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/QRCode'
        '404':
          description: Ticket not found
        '500':
          description: Internal server error
components:
  schemas:
    QRCode:
      type: object
      properties:
        id:
          type: string
          format: uuid
          readOnly: true
        image:
          type: string
          description: Base64-encoded QR code image
        client_number:
          type: string
          description: Client phone number
        used:
          type: boolean
          description: Whether the ticket has been scanned
    QRCodeInput:
      type: object
      required:
        - image
        - client_number
      properties:
        image:
          type: string
          description: Base64-encoded QR code image
        client_number:
          type: string
          description: Client phone number
`

// swaggerUI serves an HTML page that renders the OpenAPI spec via Swagger UI.
func swaggerUI(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/docs/openapi.yaml" {
		w.Header().Set("Content-Type", "application/yaml")
		w.Write([]byte(openAPISpec))
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8"/>
  <title>Tickets API — Docs</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist/swagger-ui.css"/>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"></script>
  <script>
    SwaggerUIBundle({
      url: "/docs/openapi.yaml",
      dom_id: "#swagger-ui",
      presets: [SwaggerUIBundle.presets.apis, SwaggerUIBundle.SwaggerUIStandalonePreset],
      layout: "BaseLayout"
    });
  </script>
</body>
</html>`))
}
