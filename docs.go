package main

import (
	"net/http"
)

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
          description: Invalid request body or validation error
        '409':
          description: Client number already exists
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
    delete:
      summary: Delete a ticket
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '204':
          description: Ticket deleted
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
        '409':
          description: Ticket already used
        '500':
          description: Internal server error
  /image/{id}:
    get:
      summary: Get QR code image
      description: Returns a PNG image of the QR code.
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '200':
          description: QR code image
          content:
            image/png:
              schema:
                type: string
                format: binary
        '404':
          description: Ticket not found
  /scan/{id}:
    get:
      summary: Scan a ticket
      description: Marks the ticket as used and returns a styled HTML page with the scan result.
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '200':
          description: Scan result page
          content:
            text/html:
              schema:
                type: string
        '404':
          description: Ticket not found
        '409':
          description: Ticket already used
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
          description: URL to the QR code image
        client_number:
          type: string
          description: Client phone number
        used:
          type: boolean
          description: Whether the ticket has been scanned
        created_at:
          type: string
          format: date-time
          description: Timestamp when the ticket was created
        used_at:
          type: string
          format: date-time
          nullable: true
          description: Timestamp when the ticket was scanned (null if not used)
    QRCodeInput:
      type: object
      required:
        - client_number
      properties:
        client_number:
          type: string
          description: Client phone number (must be numeric)
`

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
