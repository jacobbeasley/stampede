import '../css/application.css'
import HelloWorld from './components/HelloWorld.svelte'
import TodoApp from './components/TodoApp.svelte'
import { mount } from 'svelte'

const appTarget = document.getElementById('app')
if (appTarget) {
  mount(HelloWorld, { target: appTarget })
}

const todoTarget = document.getElementById('todo-app')
if (todoTarget) {
  mount(TodoApp, { target: todoTarget })
}
