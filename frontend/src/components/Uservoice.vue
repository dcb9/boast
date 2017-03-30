<script>
if (process.env.NODE_ENV == "production") {
  // UserVoice JavaScript SDK (only needed once on a page)
  (function(){var uv=document.createElement('script');uv.type='text/javascript';uv.async=true;uv.src='//widget.uservoice.com/7ZYFHPX8OMIzjhVqREeXA.js';var s=document.getElementsByTagName('script')[0];s.parentNode.insertBefore(uv,s)})()

  UserVoice = window.UserVoice || [];
  UserVoice.push(['showTab', 'classic_widget', {
    mode: 'full',
    primary_color: '#fa8c28',
    link_color: '#0a8cc6',
    default_mode: 'support',
    forum_id: 596554,
    tab_label: 'Feedback & Support',
    tab_color: '#fa8c28',
    tab_position: 'middle-right',
    tab_inverted: false
  }]);
}
</script>
