<script>
  let activeTab = $state("tab1");
</script>

<div class="mt-12 w-full space-y-12">
  <div class="flex items-center justify-between">
    <div class="text-2xl font-bold">DaisyUI Kitchen Sink</div>

    <div class="dropdown dropdown-end">
      <div tabindex="0" role="button" class="btn btn-outline m-1">
        Theme
        <svg
          width="12px"
          height="12px"
          class="h-2 w-2 fill-current opacity-60 inline-block ml-1"
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 2048 2048"
          ><path d="M1799 349l242 241-1017 1017L7 590l242-241 775 775 775-775z"
          ></path></svg
        >
      </div>
      <ul
        role="menu"
        tabindex="0"
        class="dropdown-content bg-base-300 rounded-box z-[1] w-52 p-2 shadow-2xl h-[28rem] flex-nowrap overflow-y-auto"
      >
        <li>
          <input
            type="radio"
            name="theme-dropdown"
            class="theme-controller btn btn-sm btn-block btn-ghost justify-start"
            aria-label="Light"
            value="light"
          />
        </li>
        <li>
          <input
            type="radio"
            name="theme-dropdown"
            class="theme-controller btn btn-sm btn-block btn-ghost justify-start"
            aria-label="Dark"
            value="dark"
          />
        </li>
        <li>
          <input
            type="radio"
            name="theme-dropdown"
            class="theme-controller btn btn-sm btn-block btn-ghost justify-start"
            aria-label="Cupcake"
            value="cupcake"
          />
        </li>
        <li>
          <input
            type="radio"
            name="theme-dropdown"
            class="theme-controller btn btn-sm btn-block btn-ghost justify-start"
            aria-label="Bumblebee"
            value="bumblebee"
          />
        </li>
        <li>
          <input
            type="radio"
            name="theme-dropdown"
            class="theme-controller btn btn-sm btn-block btn-ghost justify-start"
            aria-label="Emerald"
            value="emerald"
          />
        </li>
        <li>
          <input
            type="radio"
            name="theme-dropdown"
            class="theme-controller btn btn-sm btn-block btn-ghost justify-start"
            aria-label="Corporate"
            value="corporate"
          />
        </li>
        <li>
          <input
            type="radio"
            name="theme-dropdown"
            class="theme-controller btn btn-sm btn-block btn-ghost justify-start"
            aria-label="Synthwave"
            value="synthwave"
          />
        </li>
        <li>
          <input
            type="radio"
            name="theme-dropdown"
            class="theme-controller btn btn-sm btn-block btn-ghost justify-start"
            aria-label="Retro"
            value="retro"
          />
        </li>
        <li>
          <input
            type="radio"
            name="theme-dropdown"
            class="theme-controller btn btn-sm btn-block btn-ghost justify-start"
            aria-label="Cyberpunk"
            value="cyberpunk"
          />
        </li>
        <li>
          <input
            type="radio"
            name="theme-dropdown"
            class="theme-controller btn btn-sm btn-block btn-ghost justify-start"
            aria-label="Valentine"
            value="valentine"
          />
        </li>
        <li>
          <input
            type="radio"
            name="theme-dropdown"
            class="theme-controller btn btn-sm btn-block btn-ghost justify-start"
            aria-label="Halloween"
            value="halloween"
          />
        </li>
        <li>
          <input
            type="radio"
            name="theme-dropdown"
            class="theme-controller btn btn-sm btn-block btn-ghost justify-start"
            aria-label="Garden"
            value="garden"
          />
        </li>
        <li>
          <input
            type="radio"
            name="theme-dropdown"
            class="theme-controller btn btn-sm btn-block btn-ghost justify-start"
            aria-label="Forest"
            value="forest"
          />
        </li>
        <li>
          <input
            type="radio"
            name="theme-dropdown"
            class="theme-controller btn btn-sm btn-block btn-ghost justify-start"
            aria-label="Aqua"
            value="aqua"
          />
        </li>
        <li>
          <input
            type="radio"
            name="theme-dropdown"
            class="theme-controller btn btn-sm btn-block btn-ghost justify-start"
            aria-label="Lofi"
            value="lofi"
          />
        </li>
        <li>
          <input
            type="radio"
            name="theme-dropdown"
            class="theme-controller btn btn-sm btn-block btn-ghost justify-start"
            aria-label="Pastel"
            value="pastel"
          />
        </li>
        <li>
          <input
            type="radio"
            name="theme-dropdown"
            class="theme-controller btn btn-sm btn-block btn-ghost justify-start"
            aria-label="Fantasy"
            value="fantasy"
          />
        </li>
        <li>
          <input
            type="radio"
            name="theme-dropdown"
            class="theme-controller btn btn-sm btn-block btn-ghost justify-start"
            aria-label="Wireframe"
            value="wireframe"
          />
        </li>
        <li>
          <input
            type="radio"
            name="theme-dropdown"
            class="theme-controller btn btn-sm btn-block btn-ghost justify-start"
            aria-label="Black"
            value="black"
          />
        </li>
        <li>
          <input
            type="radio"
            name="theme-dropdown"
            class="theme-controller btn btn-sm btn-block btn-ghost justify-start"
            aria-label="Luxury"
            value="luxury"
          />
        </li>
        <li>
          <input
            type="radio"
            name="theme-dropdown"
            class="theme-controller btn btn-sm btn-block btn-ghost justify-start"
            aria-label="Dracula"
            value="dracula"
          />
        </li>
        <li>
          <input
            type="radio"
            name="theme-dropdown"
            class="theme-controller btn btn-sm btn-block btn-ghost justify-start"
            aria-label="Cmyk"
            value="cmyk"
          />
        </li>
        <li>
          <input
            type="radio"
            name="theme-dropdown"
            class="theme-controller btn btn-sm btn-block btn-ghost justify-start"
            aria-label="Autumn"
            value="autumn"
          />
        </li>
        <li>
          <input
            type="radio"
            name="theme-dropdown"
            class="theme-controller btn btn-sm btn-block btn-ghost justify-start"
            aria-label="Business"
            value="business"
          />
        </li>
        <li>
          <input
            type="radio"
            name="theme-dropdown"
            class="theme-controller btn btn-sm btn-block btn-ghost justify-start"
            aria-label="Acid"
            value="acid"
          />
        </li>
        <li>
          <input
            type="radio"
            name="theme-dropdown"
            class="theme-controller btn btn-sm btn-block btn-ghost justify-start"
            aria-label="Lemonade"
            value="lemonade"
          />
        </li>
        <li>
          <input
            type="radio"
            name="theme-dropdown"
            class="theme-controller btn btn-sm btn-block btn-ghost justify-start"
            aria-label="Night"
            value="night"
          />
        </li>
        <li>
          <input
            type="radio"
            name="theme-dropdown"
            class="theme-controller btn btn-sm btn-block btn-ghost justify-start"
            aria-label="Coffee"
            value="coffee"
          />
        </li>
        <li>
          <input
            type="radio"
            name="theme-dropdown"
            class="theme-controller btn btn-sm btn-block btn-ghost justify-start"
            aria-label="Winter"
            value="winter"
          />
        </li>
        <li>
          <input
            type="radio"
            name="theme-dropdown"
            class="theme-controller btn btn-sm btn-block btn-ghost justify-start"
            aria-label="Dim"
            value="dim"
          />
        </li>
        <li>
          <input
            type="radio"
            name="theme-dropdown"
            class="theme-controller btn btn-sm btn-block btn-ghost justify-start"
            aria-label="Nord"
            value="nord"
          />
        </li>
        <li>
          <input
            type="radio"
            name="theme-dropdown"
            class="theme-controller btn btn-sm btn-block btn-ghost justify-start"
            aria-label="Sunset"
            value="sunset"
          />
        </li>
      </ul>
    </div>
  </div>
  <div class="divider mt-0 mb-8"></div>

  <!-- Actions -->
  <section class="space-y-4">
    <h3 class="text-xl font-bold">Actions</h3>
    <div class="flex flex-wrap gap-4 items-center">
      <button class="btn">Default</button>
      <button class="btn btn-primary">Primary</button>
      <button class="btn btn-secondary">Secondary</button>
      <button class="btn btn-accent btn-outline">Accent Outline</button>
      <button class="btn btn-info btn-square" aria-label="Info">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="h-6 w-6"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          ><path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
          /></svg
        >
      </button>

      <div class="dropdown dropdown-end">
        <div tabindex="0" role="button" class="btn m-1">Dropdown</div>
        <ul
          role="menu"
          tabindex="0"
          class="dropdown-content z-[1] menu p-2 shadow bg-base-100 rounded-box w-52"
        >
          <li><a href="#item1">Item 1</a></li>
          <li><a href="#item2">Item 2</a></li>
        </ul>
      </div>

      <button
        class="btn btn-error"
        onclick={() => document.getElementById("my_modal_1").showModal()}
        >Open Modal</button
      >
      <dialog id="my_modal_1" class="modal">
        <div class="modal-box">
          <h3 class="font-bold text-lg">Hello!</h3>
          <p class="py-4">Press ESC key or click the button below to close</p>
          <div class="modal-action">
            <form method="dialog">
              <button class="btn">Close</button>
            </form>
          </div>
        </div>
      </dialog>
    </div>
  </section>

  <!-- Data Display -->
  <section class="space-y-4">
    <h3 class="text-xl font-bold">Data Display</h3>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div class="card w-full bg-base-100 shadow-xl">
        <figure>
          <img
            src="https://img.daisyui.com/images/stock/photo-1606107557195-0e29a4b5b4aa.webp"
            alt="Shoes"
          />
        </figure>
        <div class="card-body">
          <h2 class="card-title">
            Shoes!
            <div class="badge badge-secondary">NEW</div>
          </h2>
          <p>If a dog chews shoes whose shoes does he choose?</p>
          <div class="card-actions justify-end">
            <div class="badge badge-outline">Fashion</div>
            <div class="badge badge-outline">Products</div>
          </div>
        </div>
      </div>

      <div class="stats shadow w-full">
        <div class="stat">
          <div class="stat-figure text-primary">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              class="inline-block w-8 h-8 stroke-current"
              ><path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"
              ></path></svg
            >
          </div>
          <div class="stat-title">Total Likes</div>
          <div class="stat-value text-primary">25.6K</div>
          <div class="stat-desc">21% more than last month</div>
        </div>
      </div>

      <div class="flex flex-col gap-4">
        <div class="chat chat-start">
          <div class="chat-bubble chat-bubble-primary">
            What kind of nonsense is this
          </div>
        </div>
        <div class="chat chat-end">
          <div class="chat-bubble chat-bubble-secondary">
            Put me on the Council and not make me a Master!??
          </div>
        </div>

        <div class="collapse bg-base-200">
          <input type="radio" name="my-accordion-1" checked="checked" />
          <div class="collapse-title text-xl font-medium">
            Click to open this one and close others
          </div>
          <div class="collapse-content">
            <p>hello</p>
          </div>
        </div>
      </div>
    </div>
  </section>

  <!-- Navigation -->
  <section class="space-y-4">
    <h3 class="text-xl font-bold">Navigation</h3>
    <div class="navbar bg-base-300 rounded-box">
      <div class="flex-1">
        <a href="/" class="btn btn-ghost text-xl">daisyUI</a>
      </div>
      <div class="flex-none">
        <ul class="menu menu-horizontal px-1">
          <li><a href="#link">Link</a></li>
          <li>
            <details>
              <summary>Parent</summary>
              <ul class="p-2 bg-base-100 rounded-t-none">
                <li><a href="#link1">Link 1</a></li>
                <li><a href="#link2">Link 2</a></li>
              </ul>
            </details>
          </li>
        </ul>
      </div>
    </div>

    <div role="tablist" class="tabs tabs-boxed">
      <button
        type="button"
        role="tab"
        class="tab"
        class:tab-active={activeTab === "tab1"}
        onclick={() => (activeTab = "tab1")}>Tab 1</button
      >
      <button
        type="button"
        role="tab"
        class="tab"
        class:tab-active={activeTab === "tab2"}
        onclick={() => (activeTab = "tab2")}>Tab 2</button
      >
      <button
        type="button"
        role="tab"
        class="tab"
        class:tab-active={activeTab === "tab3"}
        onclick={() => (activeTab = "tab3")}>Tab 3</button
      >
    </div>

    <ul class="steps w-full">
      <li class="step step-primary">Register</li>
      <li class="step step-primary">Choose plan</li>
      <li class="step">Purchase</li>
      <li class="step">Receive Product</li>
    </ul>
  </section>

  <!-- Data Input -->
  <section class="space-y-4">
    <h3 class="text-xl font-bold">Data Input</h3>
    <div class="flex flex-col md:flex-row gap-6">
      <div class="flex flex-col gap-4 flex-1">
        <label class="input input-bordered flex items-center gap-2">
          Email
          <input type="text" class="grow" placeholder="daisy@site.com" />
        </label>

        <select class="select select-bordered w-full">
          <option disabled selected>Who shot first?</option>
          <option>Han Solo</option>
          <option>Greedo</option>
        </select>

        <input
          type="range"
          min="0"
          max="100"
          value="40"
          class="range range-primary"
        />
      </div>

      <div class="flex flex-col gap-4 flex-1 justify-center">
        <div class="form-control">
          <label class="label cursor-pointer">
            <span class="label-text">Remember me</span>
            <input
              type="checkbox"
              checked="checked"
              class="checkbox checkbox-primary"
            />
          </label>
        </div>

        <div class="form-control">
          <label class="label cursor-pointer">
            <span class="label-text">Enable notifications</span>
            <input type="checkbox" class="toggle toggle-secondary" checked />
          </label>
        </div>

        <div class="rating">
          <input
            type="radio"
            name="rating-2"
            class="mask mask-star-2 bg-orange-400"
          />
          <input
            type="radio"
            name="rating-2"
            class="mask mask-star-2 bg-orange-400"
            checked="checked"
          />
          <input
            type="radio"
            name="rating-2"
            class="mask mask-star-2 bg-orange-400"
          />
          <input
            type="radio"
            name="rating-2"
            class="mask mask-star-2 bg-orange-400"
          />
          <input
            type="radio"
            name="rating-2"
            class="mask mask-star-2 bg-orange-400"
          />
        </div>
      </div>
    </div>
  </section>

  <!-- Feedback -->
  <section class="space-y-4 pb-12">
    <h3 class="text-xl font-bold">Feedback</h3>
    <div class="flex flex-col gap-4">
      <div role="alert" class="alert alert-success">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="stroke-current shrink-0 h-6 w-6"
          fill="none"
          viewBox="0 0 24 24"
          ><path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
          /></svg
        >
        <span>Your purchase has been confirmed!</span>
      </div>

      <div class="flex items-center gap-6 justify-center p-4">
        <span class="loading loading-spinner loading-lg text-primary"></span>
        <span class="loading loading-dots loading-lg text-secondary"></span>
        <span class="loading loading-ring loading-lg text-accent"></span>

        <div
          class="radial-progress text-primary"
          style="--value:70; --size:3rem;"
          role="progressbar"
        >
          70%
        </div>
      </div>
    </div>
  </section>
</div>
