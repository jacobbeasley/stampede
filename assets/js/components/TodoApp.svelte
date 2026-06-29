<script>
  import { onMount } from "svelte";

  let todos = $state([]);
  let loading = $state(true);
  let error = $state("");
  let newTodoTitle = $state("");

  // Drag and drop state
  let dragId = $state(null);
  let dragOverId = $state(null);
  let dragDirection = $state(null); // 'up' or 'down'

  // Helper for parsing csrf token from meta tag
  function getCsrfToken() {
    const meta = document.querySelector('meta[name="csrf-token"]');
    return meta ? meta.getAttribute("content") : "";
  }

  async function loadTodos() {
    loading = true;
    error = "";
    try {
      const res = await fetch("/api/todos/");
      if (res.ok) {
        let text = await res.text();
        if (text) {
          todos = JSON.parse(text);
          if (!Array.isArray(todos)) {
            todos = [];
          }
        } else {
            todos = [];
        }
      } else {
        error = "Failed to load to-dos.";
      }
    } catch (err) {
      error = "An error occurred while loading to-dos.";
      console.error(err);
    } finally {
      loading = false;
    }
  }

  async function addTodo() {
    if (!newTodoTitle.trim()) return;

    try {
      const res = await fetch("/api/todos/", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "X-CSRF-Token": getCsrfToken(),
        },
        body: JSON.stringify({
          title: newTodoTitle,
          description: "", // Currently just setting empty description
        }),
      });

      if (res.ok) {
        const todo = await res.json();
        todos = [...todos, todo];
        newTodoTitle = "";
      } else {
        error = "Failed to add to-do.";
      }
    } catch (err) {
      error = "An error occurred while adding to-do.";
      console.error(err);
    }
  }

  async function toggleComplete(todo) {
    const updatedTodo = { ...todo, is_completed: !todo.is_completed };

    // Optimistic update
    todos = todos.map((t) => (t.id === todo.id ? updatedTodo : t));

    try {
      const res = await fetch(`/api/todos/${todo.id}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          "X-CSRF-Token": getCsrfToken(),
        },
        body: JSON.stringify(updatedTodo),
      });

      if (!res.ok) {
        // Revert on failure
        todos = todos.map((t) => (t.id === todo.id ? todo : t));
        error = "Failed to update to-do.";
      }
    } catch (err) {
      // Revert on failure
      todos = todos.map((t) => (t.id === todo.id ? todo : t));
      error = "An error occurred while updating to-do.";
      console.error(err);
    }
  }

  async function updateTitle(todo, newTitle) {
      if (!newTitle.trim()) return;
      if (todo.title === newTitle) return;

      const updatedTodo = { ...todo, title: newTitle };
      todos = todos.map((t) => (t.id === todo.id ? updatedTodo : t));

      try {
        const res = await fetch(`/api/todos/${todo.id}`, {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
            "X-CSRF-Token": getCsrfToken(),
          },
          body: JSON.stringify(updatedTodo),
        });

        if (!res.ok) {
            todos = todos.map((t) => (t.id === todo.id ? todo : t));
            error = "Failed to update to-do title.";
        }
      } catch(err) {
          todos = todos.map((t) => (t.id === todo.id ? todo : t));
          error = "An error occurred while updating to-do title.";
          console.error(err);
      }
  }

  async function deleteTodo(id) {
    // Optimistic update
    const previousTodos = [...todos];
    todos = todos.filter((t) => t.id !== id);

    try {
      const res = await fetch(`/api/todos/${id}`, {
        method: "DELETE",
        headers: {
          "X-CSRF-Token": getCsrfToken(),
        },
      });

      if (!res.ok) {
        // Revert on failure
        todos = previousTodos;
        error = "Failed to delete to-do.";
      }
    } catch (err) {
      // Revert on failure
      todos = previousTodos;
      error = "An error occurred while deleting to-do.";
      console.error(err);
    }
  }

  // --- Drag and Drop Handlers ---

  function handleDragStart(e, id) {
    dragId = id;
    e.dataTransfer.effectAllowed = "move";
    e.dataTransfer.setData("text/plain", id);
    // Slight delay to allow the drag image to be captured before we fade the original element
    setTimeout(() => {
        if (e.target) {
            e.target.classList.add("opacity-50");
        }
    }, 0);
  }

  function handleDragEnd(e) {
    if (e.target) {
        e.target.classList.remove("opacity-50");
    }
    dragId = null;
    dragOverId = null;
    dragDirection = null;
  }

  function handleDragOver(e, id) {
    e.preventDefault(); // Necessary to allow dropping
    if (id === dragId) return;

    dragOverId = id;

    // Determine drop placement relative to the target item (above or below)
    const targetElement = e.currentTarget;
    if (targetElement) {
        const bounding = targetElement.getBoundingClientRect();
        const offset = bounding.y + bounding.height / 2;
        if (e.clientY - offset > 0) {
            dragDirection = 'down';
        } else {
            dragDirection = 'up';
        }
    }
  }

  function handleDragLeave(e, id) {
     if (dragOverId === id) {
         dragOverId = null;
         dragDirection = null;
     }
  }

  async function handleDrop(e, targetId) {
    e.preventDefault();
    if (dragId === targetId || !dragId) return;

    const sourceIndex = todos.findIndex((t) => t.id === dragId);
    let targetIndex = todos.findIndex((t) => t.id === targetId);

    if (sourceIndex === -1 || targetIndex === -1) return;

    // Adjust target index based on drag direction
    if (dragDirection === 'down' && sourceIndex < targetIndex) {
        // We're moving an item down, and dropping below an item. Target index is fine.
    } else if (dragDirection === 'down' && sourceIndex > targetIndex) {
         // Moving item up, dropping below an item. Insert after.
         targetIndex++;
    } else if (dragDirection === 'up' && sourceIndex < targetIndex) {
         // Moving item down, dropping above an item. Insert before.
         targetIndex--;
    } else if (dragDirection === 'up' && sourceIndex > targetIndex) {
         // Moving item up, dropping above an item. Target index is fine.
    }

    const previousTodos = [...todos];
    const newTodos = [...todos];

    // Remove the dragged item
    const [draggedItem] = newTodos.splice(sourceIndex, 1);

    // Insert it at the new position
    newTodos.splice(targetIndex, 0, draggedItem);

    todos = newTodos;

    dragId = null;
    dragOverId = null;
    dragDirection = null;

    // Send updated order to server
    const orderIds = newTodos.map((t) => t.id);
    try {
        const res = await fetch("/api/todos/reorder", {
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
                "X-CSRF-Token": getCsrfToken(),
            },
            body: JSON.stringify({ order: orderIds }),
        });

        if (!res.ok) {
             // Revert on failure
             todos = previousTodos;
             error = "Failed to save new order.";
        }
    } catch(err) {
        // Revert on failure
        todos = previousTodos;
        error = "An error occurred while saving new order.";
        console.error(err);
    }
  }

  onMount(() => {
    loadTodos();
  });
</script>

<div class="card bg-base-100 shadow-xl w-full max-w-3xl mx-auto mt-8">
  <div class="card-body">
    <h2 class="card-title text-2xl mb-4 font-bold">My Tasks</h2>

    {#if error}
      <div class="alert alert-error mb-4">
        <span>{error}</span>
        <button class="btn btn-ghost btn-sm" onclick={() => (error = "")}>✕</button>
      </div>
    {/if}

    <form
      class="flex gap-2 mb-6"
      onsubmit={(e) => {
        e.preventDefault();
        addTodo();
      }}
    >
      <input
        type="text"
        placeholder="Add a new task..."
        class="input input-bordered flex-grow"
        bind:value={newTodoTitle}
      />
      <button type="submit" class="btn btn-primary" disabled={!newTodoTitle.trim()}>
        Add
      </button>
    </form>

    {#if loading}
      <div class="flex justify-center p-8">
        <span class="loading loading-spinner loading-lg text-primary"></span>
      </div>
    {:else if todos.length === 0}
      <div class="text-center py-8 text-base-content/50">
        <p>No tasks yet. Add one above!</p>
      </div>
    {:else}
      <ul class="space-y-2">
        {#each todos as todo (todo.id)}
          <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
          <li
            class="flex items-center gap-3 p-3 bg-base-200 rounded-lg group cursor-move transition-colors
                   {dragOverId === todo.id && dragDirection === 'up' ? 'border-t-2 border-primary' : ''}
                   {dragOverId === todo.id && dragDirection === 'down' ? 'border-b-2 border-primary' : ''}
                   {dragId === todo.id ? 'opacity-50' : ''}"
            draggable="true"
            ondragstart={(e) => handleDragStart(e, todo.id)}
            ondragend={handleDragEnd}
            ondragover={(e) => handleDragOver(e, todo.id)}
            ondragleave={(e) => handleDragLeave(e, todo.id)}
            ondrop={(e) => handleDrop(e, todo.id)}
          >
            <!-- Drag handle icon (optional visual cue) -->
            <div class="text-base-content/30 cursor-move" aria-hidden="true">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 8h16M4 16h16" />
                </svg>
            </div>

            <input
              type="checkbox"
              class="checkbox checkbox-primary cursor-pointer"
              checked={todo.is_completed}
              onchange={() => toggleComplete(todo)}
            />
            <input
              type="text"
              class="input input-ghost flex-grow focus:bg-base-100 focus:outline-none focus:ring-2 focus:ring-primary h-auto py-1 px-2 cursor-text {todo.is_completed ? 'line-through text-base-content/50' : ''}"
              value={todo.title}
              onblur={(e) => updateTitle(todo, e.target.value)}
              onkeydown={(e) => {
                  if (e.key === 'Enter') {
                      e.target.blur();
                  }
              }}
            />
            <button
              class="btn btn-ghost btn-sm btn-square text-error opacity-0 group-hover:opacity-100 transition-opacity focus:opacity-100"
              aria-label="Delete task"
              onclick={() => deleteTodo(todo.id)}
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
            </button>
          </li>
        {/each}
      </ul>
    {/if}
  </div>
</div>
