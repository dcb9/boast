<template>
  <div id="app">
    <b-navbar type="inverse" class="navbar-fixed-top bg-inverse">
      <a class="navbar-brand" href="#">Boast</a>
    </b-navbar>
    <router-view></router-view>
  </div>
</template>

<script>

export default {
  name: 'app'
}
</script>

<style>
@import "./assets/css/bootstrap.min.css";
@import "./assets/css/font-awesome.min.css";

table.params {
    font-size: 12px; font-family: Courier, monospace;
}

.transaction-list tr {
    cursor: pointer;
}

.transaction-list  tr:hover {
    background-color: #ddd;
}

tr.selected, tr.selected:hover {
    background-color: #ff9999;
    background-color: #000000;
    color:white;
}

.path {
    width: 300px;
}

.wrapped {
    word-wrap: break-word;
    overflow: hidden;
}
</style>
