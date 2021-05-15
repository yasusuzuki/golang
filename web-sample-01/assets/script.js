
//"契約DB詳細"ボタンや”事故DB詳細"ボタンを押下時に、遷移先をボタンに応じて動的に変更する機能
//前提：　<FORM>のNAMEが”DFORM"でなければならない.
function setActionToDataForm(func) {
  document.DFORM.action = func;
}
//ヘッダーのチェックボックスを押すと、すべての行のチェックと連動する
function toggleAllMsg(fld, prefix) {
  var len = prefix.length;
  for (i = 0; i < document.DFORM.elements.length; i++) {
    if ((document.DFORM.elements[i].type == "checkbox") && (document.DFORM.elements[i].name == prefix))
      document.DFORM.elements[i].checked = fld.checked;
  }
}
