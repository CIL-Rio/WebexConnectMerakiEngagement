// Add an event listener to the 'phoneNumber' input field
document.getElementById("phoneNumber").addEventListener("input", function (e) {
  // Get the target element, cursor position, and length of the input
  var target = e.target,
    position = target.selectionEnd,
    length = target.value.length;

  // Remove all non-digit characters from the input
  target.value = target.value.replace(/\D/g, "");

  // Limit the input to 13 digits
  if (target.value.length > 13) {
    target.value = target.value.slice(0, 13);
  }

  // Set the hidden input value to the formatted phone number
  document.getElementById("formattedPhoneNumber").value = target.value;

  // Format the phone number
  target.value = target.value.replace(/^(\d{2})(\d)/g, "$1 $2");
  target.value = target.value.replace(/(\d{2}) (\d{2})(\d)/g, "$1 $2 $3");
  target.value = target.value.replace(/(\d)(\d{4})$/, "$1-$2");

  // Adjust the cursor position
  if (position !== length) {
    if (e.data && /\D/g.test(e.data)) {
      target.selectionEnd = position - 1;
    } else {
      target.selectionEnd = position;
    }
  }
});

// Get the submit button
var submitButton = document.getElementById("submitButton");

// Disable or enable the submit button
submitButton.disabled = disableSubmit;

// If the submit button is disabled, start a countdown
if (disableSubmit) {
  submitButton.textContent = "Reenviar mensagem";
  var countdownInterval = setInterval(function () {
    countdown--;
    if (countdown <= 0) {
      clearInterval(countdownInterval);
      submitButton.value = "Reenviar mensagem";
      submitButton.disabled = false;
    } else {
      submitButton.value = "Reenviar mensagem (" + countdown + ")";
    }
  }, 1000);
}

// Trigger the 'input' event on the 'phoneNumber' input field when the page loads
window.onload = function () {
  var event = new Event("input");
  document.getElementById("phoneNumber").dispatchEvent(event);
};

if (messageSuccess) {
  Swal.fire({
    icon: "success",
    title: "Sucesso",
    text: "Enviamos uma mensagem com link de verificação para seu WhatsApp. Por favor, verifique a mensagem para continuar.",
  });
} else if (message) {
  document.getElementById("errorMessageText").textContent = message;
  document.getElementById("errorMessage").style.display = "block";
}
