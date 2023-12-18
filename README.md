# Zopsmart Project

This is a simple project on React and GoFr framework for Zopsmart Assignment Submission.

## API DESCRIPTION
# Todo Model Endpoints:

- Create Todo (POST):

    - Endpoint: /todos
    - Description: Create a new todo item.
    - Request Body:
    - JSON
    
    {
        "title": "Buy groceries",
        "description": "Get milk, eggs, and bread"
    }
    - Response: Returns the newly created todo with an assigned ID.
- Get Todo by ID (GET):

    - Endpoint: /todos/{id}
    - Description: Retrieve todo details by its unique ID.
    - Response: Returns todo information.
- Update Todo (PUT/PATCH):

    - Endpoint: /todos/{id}
    - Description: Update todo details.
    - Request Body (example for marking as completed):
    - JSON
  
    {
        "completed": true
    }
    - Response: Returns the updated todo.
- Delete Todo (DELETE):

    - Endpoint: /todos/{id}
    - Description: Delete a todo by its ID.
    - Response: Returns a success message.
 
# User Endpoints
    - /register
    - /login
      
 
