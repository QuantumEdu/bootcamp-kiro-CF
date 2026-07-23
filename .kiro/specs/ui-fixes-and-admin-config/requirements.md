# Requirements Document

## Introduction

This spec covers a set of UI improvements and a new admin configuration feature for the POS AI-First application. It includes: adding a visible logout button to the sidebar, connecting the existing product creation form in the UI, implementing client (clientes) management CRUD, ensuring HTMX fragments refresh correctly on navigation, and providing an admin-only settings page to securely store and manage an API key.

## Glossary

- **System**: The POS AI-First web application backend (Go + chi + HTMX)
- **Sidebar**: The left navigation panel rendered by layout.html, visible to all authenticated users
- **Authenticated_User**: A user with a valid session (user_id present in session store)
- **Admin_User**: An Authenticated_User whose session role is "admin"
- **Products_Page**: The page served at GET /productos displaying the product list
- **Product_Form**: The HTML form template at templates/products/form.html for creating a new product
- **Clients_Page**: The page listing all clients, served at GET /clientes
- **Client_Form**: The HTML form for creating a new client, served at GET /clientes/new
- **HTMX_Fragment**: A partial HTML response loaded via HTMX into the page without full reload
- **Admin_Config_Page**: An admin-only settings page served at GET /admin/config
- **Configuracion_Table**: The SQLite table `configuracion` storing key-value pairs for system settings
- **API_Key**: The OpenRouter API key stored encrypted in the Configuracion_Table
- **SESSION_SECRET**: The environment variable used to derive the AES-GCM encryption key

## Requirements

### Requirement 1: Logout Button in Sidebar

**User Story:** As an authenticated user, I want a logout button visible in the sidebar footer, so that I can end my session from any page.

#### Acceptance Criteria

1. THE Sidebar SHALL display a logout button in the footer section for every Authenticated_User.
2. WHEN the Authenticated_User clicks the logout button, THE System SHALL send a POST request to /logout.
3. WHEN the POST /logout request is processed, THE System SHALL destroy the user session and redirect to /login.
4. THE Sidebar SHALL render the logout button below the user info section in the footer area.

### Requirement 2: Product Creation UI Button

**User Story:** As an authenticated user, I want a "Nuevo Producto" button on the products page, so that I can navigate to the product creation form.

#### Acceptance Criteria

1. THE Products_Page SHALL display a "Nuevo Producto" button visible to all Authenticated_Users.
2. WHEN the Authenticated_User clicks the "Nuevo Producto" button, THE System SHALL navigate to GET /productos/new.
3. WHEN the GET /productos/new route is requested, THE System SHALL render the Product_Form template.

### Requirement 3: Client Listing

**User Story:** As an authenticated user, I want to view a list of all clients, so that I can see customer information.

#### Acceptance Criteria

1. THE System SHALL provide a GET /clientes route that renders the Clients_Page.
2. THE Clients_Page SHALL display all clients from the clientes table including nombre, telefono, and direccion fields.
3. THE Sidebar SHALL include a navigation link to /clientes labeled "Clientes".
4. WHEN no clients exist in the database, THE Clients_Page SHALL display an empty state message indicating no clients are registered.

### Requirement 4: Client Creation

**User Story:** As an authenticated user, I want to create a new client, so that I can register customer information for sales.

#### Acceptance Criteria

1. THE System SHALL provide a GET /clientes/new route that renders the Client_Form.
2. THE Client_Form SHALL include fields for nombre (required), telefono (optional), and direccion (optional).
3. WHEN the Authenticated_User submits the Client_Form with a valid nombre, THE System SHALL insert a new record into the clientes table and redirect to the Clients_Page.
4. IF the Authenticated_User submits the Client_Form with an empty nombre field, THEN THE System SHALL display a validation error message and not persist the record.
5. THE Clients_Page SHALL display a "Nuevo Cliente" button that navigates to GET /clientes/new.

### Requirement 5: HTMX Fragment Refresh on Navigation

**User Story:** As an authenticated user, I want page data to refresh when I navigate between sections, so that I always see current information.

#### Acceptance Criteria

1. THE System SHALL include hx-trigger="load" on HTMX fragments that load dynamic content, ensuring data is fetched on each page navigation.
2. WHEN the application server restarts, THE System SHALL serve fresh content on the next request without requiring the user to hard-refresh the browser.
3. THE System SHALL set appropriate HTTP cache-control headers on HTMX fragment responses to prevent stale cached state.

### Requirement 6: Admin API Key Configuration Page

**User Story:** As an admin user, I want a settings page where I can view and update the API key, so that I can manage external service credentials securely.

#### Acceptance Criteria

1. THE System SHALL provide a GET /admin/config route accessible only to Admin_Users.
2. IF a non-admin Authenticated_User requests GET /admin/config, THEN THE System SHALL respond with a 403 Forbidden status or redirect to the dashboard.
3. THE Admin_Config_Page SHALL display the stored API key in masked format showing only the last four characters (e.g., "****xxxx").
4. WHEN no API key is stored in the Configuracion_Table, THE Admin_Config_Page SHALL display an empty state prompting the Admin_User to configure the key.
5. THE Admin_Config_Page SHALL provide a form to submit a new or updated API key value.
6. WHEN the Admin_User submits a valid API key, THE System SHALL encrypt the key using AES-GCM with a key derived from SESSION_SECRET and store the encrypted value in the Configuracion_Table with clave "openrouter_api_key".
7. WHEN the System reads the stored API key for operational use, THE System SHALL decrypt it using the same AES-GCM key derived from SESSION_SECRET.
8. IF the SESSION_SECRET environment variable is not set or is empty, THEN THE System SHALL refuse to start and log an error indicating the secret is required.
9. THE Sidebar SHALL include a navigation link to /admin/config labeled "Configuración" visible only to Admin_Users.
