To-Do List Front-End: How It Works
This HTML file creates a fully interactive to-do list that runs in your browser and saves your tasks using localStorage.
1. HTML Structure (index.html or similar):
* <head>:
* Sets up the page title and viewport for responsiveness.
* Includes Tailwind CSS from a CDN for styling. This allows us to use utility classes directly in the HTML for a modern look without writing separate CSS files for most things.
* Includes Font Awesome for icons (like the plus and trash can icons).
* A <style> block is included for any minor custom styles, like the completed class for strikethrough text and a simple fade-in animation for new tasks. The Inter font is also suggested.
* <body>:
* Styled with a gradient background (bg-gradient-to-br from-slate-900 to-slate-700) and centers the content.
* Main Container (<div class="bg-slate-800 ...">): A styled card that holds all the to-do list elements.
* Title (<h1>): "My To-Do List".
* Input Form (<form id="taskForm">):
* An <input type="text" id="taskInput"> for typing new tasks.
* A <button type="submit"> to add the task.
* Task List (<ul id="taskList">): An unordered list where the tasks will be dynamically inserted by JavaScript.
* Empty State Message (<p id="emptyMessage">): Shown when there are no tasks. Initially hidden.
* Clear All Button (<button id="clearAllButton">): Allows deleting all tasks. Initially hidden.
* <script> (JavaScript Logic):
* All the JavaScript code is placed within <script> tags at the end of the <body> so that it runs after the HTML elements are loaded.
2. CSS (Tailwind CSS & Custom):
* Tailwind CSS: Most of the styling is done using Tailwind's utility classes directly in the HTML (e.g., bg-sky-500, p-3, rounded-lg, flex, items-center). This makes development faster and keeps styles co-located with their elements.
* Custom CSS (in <style> tags):
* .completed: Applies a line-through and lighter color to completed tasks.
* @keyframes fadeIn: A simple animation for new tasks appearing.
* .task-item: Applies the fade-in animation.
* .btn and .icon-btn: Reusable button styles for consistency and better touch targets.
3. JavaScript Logic:
* **DOM Element Selection:**
    * Variables like `taskForm`, `taskInput`, `taskList`, `emptyMessage`, and `clearAllButton` are created to reference the corresponding HTML elements using `document.getElementById()`.

* **Application State (`tasks` array):**
    * `let tasks = [];` : This array will hold all our task objects. Each task object will have an `id` (a unique number, we use `Date.now()`), `text` (the task description), and `completed` (a boolean: `true` or `false`).

* **Event Listeners:**
    * **`taskForm.addEventListener('submit', ...)`:**
        * When the form is submitted (either by clicking "Add" or pressing Enter in the input field):
            * `event.preventDefault();` stops the default browser action of reloading the page.
            * It gets the `taskText` from the input field and trims any leading/trailing whitespace.
            * If the `taskText` is not empty, it calls `addTask(taskText)`.
            * Clears the input field and sets focus back to it.
    * **`taskList.addEventListener('click', ...)`:**
        * This is an example of **event delegation**. Instead of adding an event listener to every single task item (which can be inefficient if there are many tasks), we add one listener to the parent `<ul>`.
        * When a click occurs inside the `taskList`:
            * `event.target` tells us exactly which element was clicked.
            * `target.closest('.delete-button')`: Checks if the clicked element (or one of its parents) is a delete button. If so, it gets the `taskId` from the `li` element's `dataset.id` and calls `deleteTask(taskId)`.
            * `target.type === 'checkbox' ...`: Checks if a task's checkbox was clicked. If so, it gets the `taskId` and calls `toggleComplete(taskId)`.
    * **`clearAllButton.addEventListener('click', ...)`:**
        * When the "Clear All Tasks" button is clicked:
            * It shows a confirmation dialog (`confirm(...)`).
            * If the user confirms, it clears the `tasks` array, saves the empty array to `localStorage`, and re-renders the UI.

* **Core Functions:**
    * **`addTask(text)`:**
        * Creates a new task object with a unique `id` (using `Date.now()`), the provided `text`, and `completed: false`.
        * Adds this new task object to the `tasks` array.
        * Calls `saveTasks()` to update `localStorage`.
        * Calls `renderTasks()` to update the display.
    * **`deleteTask(taskId)`:**
        * Filters the `tasks` array to remove the task with the matching `taskId`.
        * Calls `saveTasks()` and `renderTasks()`.
    * **`toggleComplete(taskId)`:**
        * Maps over the `tasks` array. If a task's `id` matches `taskId`, it creates a new task object with the `completed` status flipped (`!task.completed`).
        * Calls `saveTasks()` and `renderTasks()`.
    * **`renderTasks()`:**
        * This is the function responsible for displaying the tasks in the HTML.
        * `taskList.innerHTML = '';` clears any existing tasks from the list (so we don't get duplicates when re-rendering).
        * **Empty State:** If `tasks.length === 0`, it shows the `emptyMessage` and hides the `clearAllButton`.
        * **Displaying Tasks:** Otherwise, it hides the `emptyMessage` and shows the `clearAllButton`. It then iterates over the `tasks` array using `forEach()`:
            * For each `task` object, it dynamically creates HTML elements:
                * An `<li>` (list item) with class `task-item` and a `data-id` attribute to store the task's ID.
                * An `<input type="checkbox">` for marking the task as complete. Its `checked` state is set based on `task.completed`.
                * A `<span>` to display the `task.text`. If `task.completed` is true, the `completed` class is added for strikethrough.
                * A `<button>` with a trash icon for deleting the task.
            * These elements are appended to the `taskItem`, and then the `taskItem` is appended to the `taskList` (the `<ul>`).

* **Local Storage Functions:**
    * **`saveTasks()`:**
        * `localStorage.setItem('todoTasks', JSON.stringify(tasks));`
        * Saves the current `tasks` array into the browser's `localStorage`. Since `localStorage` can only store strings, `JSON.stringify()` is used to convert the array of objects into a JSON string. The key `'todoTasks'` is used to identify our data.
    * **`loadTasks()`:**
        * `const storedTasks = localStorage.getItem('todoTasks');` retrieves the tasks string from `localStorage`.
        * If `storedTasks` exists (meaning tasks were previously saved):
            * `tasks = JSON.parse(storedTasks);` converts the JSON string back into a JavaScript array of objects.
        * Finally, `renderTasks()` is called to display the loaded tasks.

* **Initial Load:**
    * `loadTasks();` is called once when the script first runs. This ensures that any tasks saved from a previous session are loaded and displayed.