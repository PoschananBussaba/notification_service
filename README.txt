[Overview]
This project is a Notification Service that allows users to schedule and send notifications via email using an API. It includes a frontend HTML form for input, a backend written in Go using Fiber, and a scheduler to send emails at the specified time.

[Database Schema (MySQL)]

-- Users Table
CREATE TABLE users (
    id CHAR(36) PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(15),
    role ENUM('sender', 'member', 'both') DEFAULT 'member',
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME
);

-- Target Groups Table
CREATE TABLE target_groups (
    id CHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME
);

-- Target Group Members Table
CREATE TABLE target_group_members (
    id CHAR(36) PRIMARY KEY,
    group_id CHAR(36) NOT NULL,
    user_id CHAR(36) NOT NULL,
    role ENUM('member', 'admin') DEFAULT 'member',
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    FOREIGN KEY (group_id) REFERENCES target_groups (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

-- Notifications Table
CREATE TABLE notifications (
    id CHAR(36) PRIMARY KEY,
    sender_id CHAR(36) NOT NULL,
    target_group_id CHAR(36) NOT NULL,
    subject VARCHAR(255) NOT NULL,
    message JSON NOT NULL,
    status ENUM('pending', 'in_progress', 'sent', 'failed') DEFAULT 'pending',
    priority ENUM('low', 'normal', 'high') DEFAULT 'normal',
    scheduled_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    FOREIGN KEY (sender_id) REFERENCES users (id),
    FOREIGN KEY (target_group_id) REFERENCES target_groups (id)
);

[API Endpoints]

1. Get Sender Emails
- Endpoint: GET /api/sender-emails
- Response:
[
    { "id": "1", "email": "sender1@example.com", "name": "Sender One", "role": "sender" }
]

2. Get Target Groups
- Endpoint: GET /api/target-groups
- Response:
[
    { "id": "101", "name": "Marketing Team" }
]

3. Create Notification
- Endpoint: POST /api/notifications
- Request Body:
{
    "sender_email": "1",
    "target_group": "101",
    "subject": "Test Notification",
    "message": "{\"title\": \"Test Title\", \"content\": \"Test Content\"}",
    "priority": "normal",
    "scheduled_at": "2025-01-20T15:00"
}

4. Send Scheduled Notifications
- Endpoint: POST /api/send-notifications
- Function: Sends scheduled email notifications

[Running the Program]
1. Create the MySQL database using the provided schema.
2. Set up .env for SMTP Server:
SMTP_EMAIL=your_email@gmail.com
SMTP_PASSWORD=your_password
3. Run the program:
go run main.go
4. Test API using Postman or cURL

