class IDE.StatusBar extends KDView

  constructor: (options = {}, data) ->

    options.cssClass = 'status-bar'

    super options, data

    @addSubView @status = new KDCustomHTMLView cssClass : 'status'

    @addSubView new KDCustomHTMLView
      partial  : '<cite></cite>'
      cssClass : 'icon help'
      click    : -> new HelpSupportModal

    @addSubView new KDCustomHTMLView
      partial  : '<cite></cite>'
      cssClass : 'icon github'
      click    : -> KD.utils.createExternalLink 'https://github.com/koding/IDE'

    @addSubView new KDCustomHTMLView
      partial  : '<cite></cite>'
      cssClass : 'icon shortcuts'
      click    : -> KD.getSingleton('appManager').tell 'IDE', 'showShortcutsView'

    @addSubView new KDCustomHTMLView
      partial  : '<cite></cite>'
      cssClass : 'icon participants'
      click    : => @emit 'ParticipantsModalRequired'



  showInformation: ->
    @status.updatePartial 'Click the plus button above to create a new panel'
