# 🚀 Portfolio Backend API

Welcome to the **Portfolio Backend API**! This project serves as the backend for managing portfolio projects, built with **Golang (Gin), GORM, PostgreSQL**, and JWT authentication. It allows an admin to manage projects with secure authentication and middleware protection.

---

## 📌 Features
- 🔐 **Authentication** (JWT-based login)
- 📂 **Project Management** (Create, Read, Update, Delete)
- 🖼️ **Image Upload** (Supports file uploads for project images)
- 🛡️ **Middleware Protection** (Ensures secure API access)
- 🏗️ **Clean Architecture** (Separation of concerns with controllers, services, and repositories)

---

## 🛠️ Installation & Setup

### 1️⃣ Clone the repository
```bash
$ git clone https://github.com/nazazulfiqi/be-porto-v3.git
$ cd be-porto-v3
```

### 2️⃣ Set up environment variables
Create a `.env` file and configure the following variables:
```env
PORT=8080
DB_HOST=localhost
DB_USER=
DB_PASSWORD=
DB_NAME=portofolio 
DB_PORT=5432

BASE_URL=http://localhost:8080
JWT_SECRET=your_secret_key
```

### 3️⃣ Install dependencies
```bash
$ go mod tidy
```

### 4️⃣ Run the server
```bash
$ go run main.go
```

Server will start at: **http://localhost:8080** 🚀

---

## 📡 API Endpoints

### 🔐 Authentication
| Method | Endpoint      | Description       |
|--------|-------------|------------------|
| POST   | `/auth/login` | User login |
| POST   | `/auth/register` | User login |

### 📂 Project Management
| Method | Endpoint      | Description       |
|--------|-------------|------------------|
| GET    | `/projects`   | Get all projects |
| GET    | `/projects/:id` | Get project by ID |
| POST   | `/projects`   | Create new project |
| PUT    | `/projects/:id` | Update project |
| DELETE | `/projects/:id` | Delete project |

For **image uploads**, use `multipart/form-data` with `picture_cover` and `pictures[]` fields.

---

## 📁 Folder Structure
```
📂 be-porto-v3
 ┣ 📂 controller  # API request handlers (handles HTTP requests and responses)
 ┣ 📂 database    # DB connection & migrations (handles database initialization and migrations)
 ┣ 📂 dto         # Data Transfer Objects (structures used to transfer data between layers)
 ┣ 📂 middleware  # JWT authentication middleware (handles authentication and authorization)
 ┣ 📂 models      # Database models (defines database table structures)
 ┣ 📂 repository  # Database queries (handles direct interaction with the database)
 ┣ 📂 routes      # Route definitions (maps HTTP endpoints to controllers)
 ┣ 📂 service     # Business logic layer (contains core application logic)
 ┣ 📂 uploads     # Stores uploaded images and files
 ┣ 📜 main.go     # Entry point (initializes the application and starts the server)
```

---

## ✨ Contributing
Feel free to fork this repository and submit pull requests! 😊

---

## 📬 Contact
- 🌍 **LinkedIn:** [Naza Zulfiqi](https://www.linkedin.com/in/nazazulfiqi)
- 📧 **Email:** zulfiqinaza@gmail.com

🚀 Happy coding! 🎉

