/* eslint-env browser */
var init = function() {
  setupAnalytics()
  window.addEventListener('load', redirectToNewHash)
}

var setupAnalytics = function() {
  var Analytics = window._analytics.init({
    plugins: [window.analyticsGA({ trackingId: 'UA-42258181-19' })]
  })
  Analytics.page()
}

// The previous API documentation used different fragments. We redirect those
// to the new links.
var redirectToNewHash = function() {
  var oldHash = OLD_HASH_REGEXP.exec(document.location.hash)
  if (oldHash === null) {
    return
  }

  document.location.hash = '#operation/' + oldHash[1]
}

var OLD_HASH_REGEXP = /^#\/default\/(.+)/

init()
