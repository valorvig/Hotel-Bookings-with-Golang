// attention contains a js object returned from Prompt() to select which function to use
let attention = Prompt(); // attention is {toast: toast, success: success} from Prompt's return

/*
// This should not  fire unless there's actually a date picker on that page, not on every page with the base layout
const elem = document.getElementById("reservation-dates");
const rangepicker = new DateRangePicker(elem, {
  // https://mymth.github.io/vanillajs-datepicker/#/options?id=format
  format: "yyy-mm-dd"
});
*/

// Example starter JavaScript for disabling form submissions if there are invalid fields
(function () {
  'use strict'

  // Fetch all the forms we want to apply custom Bootstrap validation styles to
  var forms = document.querySelectorAll('.needs-validation')

  // Loop over them and prevent submission
  Array.prototype.slice.call(forms)
    .forEach(function (form) {
      form.addEventListener('submit', function (event) {
        if (!form.checkValidity()) {
          event.preventDefault()
          event.stopPropagation()
        }

        form.classList.add('was-validated')
      }, false)
    })
})()
  
function notify(msg, msgType) {
  notie.alert({
    type: msgType, // optional, default = 4, enum: [1, 2, 3, 4, 5, 'success', 'warning', 'error', 'info', 'neutral']
    text: msg,
    // stay: Boolean, // optional, default = false
    // time: Number, // optional, default = 3, minimum = 1,
    // position: String // optional, default = 'top', enum: ['top', 'bottom']
  })
}

function notifyModal(title, text, icon, confirmButtonText) {
  Swal.fire({
    title: title,
    // text: text,
    html: text,
    icon: icon,
    confirmButtonText: confirmButtonText
  })
}

// display error message when the user comes to the page, using Notie
// try visiting "/reservation-summary" directly without making the reservation to see the result
{{with .Error}} //if you have some values in Error that's not empty
  notify("{{.}}", "error")
{{end}}

// if success
{{with .Flash}} 
  notify("{{.}}", "success")
{{end}}

{{with .Warning}} 
  notify("{{.}}", "warning")
{{end}}

// SweetAlert2
// javascript module - avoid hard codes, edit in one place, and change everywhere
function Prompt() {

  let toast = function(c) {
    const {
      // default
      msg = "",
      icon = "success",
      position = "top-end",
    } = c;

    const Toast = Swal.mixin({
      // try avoiding hard codes
      toast: true,
      title: msg,
      position: position,
      icon: icon,
      showConfirmButton: false,
      timer: 3000,
      timerProgressBar: true,
      didOpen: (toast) => {
        toast.addEventListener('mouseenter', Swal.stopTimer)
        toast.addEventListener('mouseleave', Swal.resumeTimer)
      }
    })

    Toast.fire({})
  }

  let success = function(c) {
    const {
      msg = "",
      title = "",
      footer = "",
    } = c;

    Swal.fire({
      icon: 'success',
      title: title,
      text: msg,
      footer: footer,
    })
  }

  let error = function(c) {
    const {
      msg = "",
      title = "",
      footer = "",
    } = c;

    Swal.fire({
      icon: 'error',
      title: title,
      text: msg,
      footer: footer,
    })
  }

  // use async since the sweet alert uses await in this one
  async function custom(c) {
    const {
      // allow these things to be specified
      icon = "",
      msg = "",
      title = "",
      showConfirmButton = true,
    } = c;
    
    // Multiple inputs modal
    const { value: result } = await Swal.fire({
      icon: icon,
      title: title,
      html: msg,
      backdrop: false,
      focusConfirm: false,
      showCancelButton: true,
      showConfirmButton: showConfirmButton,

      // Popup lifecycle hook. Synchronously runs before the popup is shown on screen.
      // Initialize the datepicker before the popup
      willOpen: () => {
        if (c.willOpen !== undefined) {
          c.willOpen(); // call willOpen (the same name in this case) from c
        }
      },
      preConfirm: () => {
        return [
          document.getElementById('start').value,
          document.getElementById('end').value
        ]
      },
      // Popup lifecycle hook. Asynchronously runs after the popup has been shown on screen.
      didOpen: () => {
        if (c.didOpen !== undefined) {
          c.didOpen();
        }
      }
    })

    // this allows to process code after the swal dialogue is closed (after users click the submit or OK button)
    if (result) {
      // https://sweetalert2.github.io/#handling-dismissals
      // if users didn't hit the cancel button, check if we have any actual values/results
      // if the result is not from clicking the cancel button on the window and not equl to empty, then do something
      if (result.dismiss !== Swal.DismissReason.cancel) { // !== is not equal "exactly"
        if (result.value !== "") {
          if (c.callback != undefined) { // undefined - not assigned a value
            c.callback(result);
          }
        } else {
          c.callback(false);
        }
      } else {
        c.callback(false);
      }
    }
    
  }

  // select each type of alert you want as a sub module
  // return object (a dictionary or a slice of choices)
  return {
    toast: toast,
    success: success,
    error: error,
    custom: custom,
  }
}