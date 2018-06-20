<template>
  <div class="container">
    <b-alert :show="!hasWebSocket" state="danger">
      Your browser does not support WebSockets.
    </b-alert>

    <b-alert :show="wsConnected && transactions.length == 0" state="info">
      No Transactions
    </b-alert>

    <b-alert :show="!wsConnected" state="danger">
      WebSocket connection is closed. <a href="./">Refresh</a>
    </b-alert>

    <div v-if="hasWebSocket && transactions.length > 0" class="main row">
      <div class="col-6">
        <h4>All Transactions</h4>
        <table class="table transaction-list">
          <tr v-bind:class="{selected: t.Selected}" v-on:click="makeActive(t, transactions)" v-for="t in transactions">
            <td class="wrapped"><div class="path">{{ t.Req.Method }} {{ t.Req.Path }}</div></td>
            <td>{{ t.Resp.Status }}</td>
            <td><span>{{ t.Duration }}</span></td>
          </tr>
        </table>
      </div>
      <div v-if="cts.ID" class="col-6">
        <div class="row">
          <div class="col-6">
            <i class="fa fa-calendar" aria-hidden="true"></i>
            <span>
              {{ cts.BeginAt | formatBeginAt }}
            </span>
          </div>
          <div  class="col-6">
            <i class="fa fa-user" aria-hidden="true"></i>
            <span class="clientIP">{{ cts.ClientIP }}</span>
          </div>
        </div>
        <hr />
        <div>
          <h3 class="row">Request</h3>
          <div>
            <p><a v-on:click="doDownloadCurl()" href="javascript:;">show cURL command</a></p>
            <span v-if="downloadCurl"><a target="_blank" href="https://www.getpostman.com/docs/importing_curl">Tips</a></span>
            <div v-if="downloadCurl" style="background-color: #c2c2c2; padding: 5px 5px 0; word-break: break-all"><pre><code style="white-space:normal;">{{ cts.Req.CurlCommand | base64Decode }}</code></pre></div>
          </div>
          <div>
            <pre><code>{{ cts.Req.RawText | base64Decode }}</code></pre>
          </div>
        </div>
        <hr style="margin: 0 0 20px" />
        <div>
          <h3>Response</h3>
          <div>
            <pre><code>{{ cts.Resp.RawText | base64Decode }}</code></pre>
          </div>
          <p><a :href="this.responseBodyURL" target="_blank">SHOW RESPONSE IN A NEW TAB</a></p>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  import Uservoice from './Uservoice'
  var Base64 = require('js-base64').Base64;

  var transactions = [];
  var currentTransaction = {};

  var durationFilter = function(ms) {
    if (ms < 1000) {
      return ms + " ms"
    }
    var s = ms / 1000
    if (s < 60) {
      return s + " s"
    }
    var m = s / 60
    if (m < 60) {
      return m + " min"
    }
    var h = m / 60
    return h + " h"
  }

  // copy from http://stackoverflow.com/questions/3177836/how-to-format-time-since-xxx-e-g-4-minutes-ago-similar-to-stack-exchange-site
  var timeSince = function (date) {

    var seconds = Math.floor((new Date() - date) / 1000);

    var interval = Math.floor(seconds / 31536000);

    if (interval >= 1) {
      return interval + " years";
    }
    interval = Math.floor(seconds / 2592000);
    if (interval >= 1) {
      return interval + " months";
    }
    interval = Math.floor(seconds / 86400);
    if (interval >= 1) {
      return interval + " days";
    }
    interval = Math.floor(seconds / 3600);
    if (interval >= 1) {
      return interval + " hours";
    }
    interval = Math.floor(seconds / 60);
    if (interval >= 1) {
      return interval + " minutes";
    }
    return Math.floor(seconds) + " seconds";
  }


  export default {
    name: 'main',
    components: { Uservoice },
    data () {
      return {
        hasWebSocket: window["WebSocket"],
        transactions: transactions,
        cts: currentTransaction,
        wsConnected: false,
        wsConn: {},
        downloadCurl: false,
        host: "",
      }
    },
    computed: {
      responseBodyURL: function() {
        return `//${this.host}/responses/${this.cts.ID}`
      },
    },
    created: function() {
      if (this.hasWebSocket) {
        if (typeof process.env.WS_HOST == "undefined") {
          this.host = window.location.host
        } else {
          this.host = process.env.WS_HOST
        }

        var conn = new WebSocket(`ws://${this.host}/ws`);
        conn.onopen = (evt) => {
          console.log("Connection connected")
          this.wsConnected = true
          this.wsConn = conn
        }
        conn.onclose = (evt) => {
          console.log("Connection closed")
          this.wsConnected = false
        }
        conn.onmessage = function (evt) {
          var data = evt.data
          if (data == "") {
            return
          }

          var tss = JSON.parse(evt.data)
          tss.forEach(function(value, key) {
            value.Duration = durationFilter(new Date(value.EndAt) - new Date(value.BeginAt))
            transactions.unshift(value)
          })
        }
      }
    },
    methods: {
      makeActive: function(t, transactions) {
        transactions.forEach(function(item, index, array){
          transactions[index].Selected = false
        });

        t.Selected = true;
        this.cts = t;
        this.downloadCurl = false;
      },
      doDownloadCurl: function() {
         this.downloadCurl = !this.downloadCurl;
      }
    },
    filters: {
      base64Decode:  (value) => {
        if (!value) return
        value = value.toString()
        return Base64.decode(value)
      },
      formatBeginAt: (value) => {
        if (!value) return
        value = value.toString()
        var d = new Date(value)

        return timeSince(d) + " ago"
      },
    }
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style>
  .alert {
    border-radius: 0;
  }
  .clientIP {
    margin-left: 8px;
  }
  .container {
    padding-top: 10px;
  }
</style>
