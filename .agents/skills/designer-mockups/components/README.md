# Mockup Components

This directory contains reusable daisyUI component patterns for mockups.

## Common Component Patterns

### Page Header with Actions
```html
<div class="mb-6 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-3">
    <div>
        <h1 class="text-3xl font-bold">Page Title</h1>
        <p class="text-base-content/70">Optional subtitle or description</p>
    </div>
    <div class="join">
        <button class="join-item btn">View</button>
        <button class="join-item btn btn-active">Edit</button>
        <button class="join-item btn">Delete</button>
    </div>
</div>
```

### Card with Metadata
```html
<div class="card bg-base-100 shadow-md border border-base-300">
    <div class="card-body">
        <h2 class="card-title">Card Title</h2>
        <p>Card description or content.</p>
        <div class="card-actions justify-end">
            <button class="btn btn-primary">Action</button>
        </div>
    </div>
</div>
```

### Status Badge
```html
<span class="badge badge-primary">Active</span>
<span class="badge badge-secondary">In Progress</span>
<span class="badge badge-success">Completed</span>
<span class="badge badge-warning">Pending</span>
<span class="badge badge-error">Error</span>
```

### Data Table
```html
<div class="overflow-x-auto border border-base-300 rounded-lg">
    <table class="table w-full">
        <thead>
            <tr>
                <th>#</th>
                <th>Name</th>
                <th>Status</th>
                <th>Actions</th>
            </tr>
        </thead>
        <tbody>
            <tr>
                <td>1</td>
                <td>Item Name</td>
                <td><span class="badge badge-success">Active</span></td>
                <td>
                    <button class="btn btn-ghost btn-xs">Edit</button>
                    <button class="btn btn-ghost btn-xs">Delete</button>
                </td>
            </tr>
        </tbody>
    </table>
</div>
```

### Modal
```html
<input type="checkbox" id="my_modal" class="modal-toggle" />
<div class="modal">
    <div class="modal-box">
        <h3 class="font-bold text-lg">Modal Title</h3>
        <p class="py-4">Modal content goes here...</p>
        <div class="modal-action">
            <label for="my_modal" class="btn">Close</label>
            <button class="btn btn-primary">Save</button>
        </div>
    </div>
</div>
```

### Alert Messages
```html
<div class="alert alert-info">
    <span>Info message text.</span>
</div>
<div class="alert alert-success">
    <span>Success message text.</span>
</div>
<div class="alert alert-warning">
    <span>Warning message text.</span>
</div>
<div class="alert alert-error">
    <span>Error message text.</span>
</div>
```

### Form Input Group
```html
<div class="form-control w-full">
    <label class="label">
        <span class="label-text font-medium">Label Text</span>
    </label>
    <input type="text" placeholder="Placeholder text" class="input input-bordered w-full" />
</div>
```

### Loading State (Skeleton)
```html
<div class="w-full">
    <div class="skeleton h-4 w-1/2 mb-2"></div>
    <div class="skeleton h-4 w-full mb-1"></div>
    <div class="skeleton h-4 w-3/4"></div>
</div>
```
