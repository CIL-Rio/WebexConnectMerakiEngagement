<?php 

include 'SplashPageAuthenticator.php';

if ($_SERVER['REQUEST_METHOD'] === 'GET' && isset($_GET['phoneNumber'], $_GET['termsAccepted']) && $_GET['termsAccepted'] === 'true') {
  try {
      $authResult = processAuthentication();
  } catch (Exception $e) {
      $authResult = ['message' => 'Ocorreu um erro', 'success' => false, 'disableSubmit' => false, 'countdown' => 0];
  }
}

?>
<!DOCTYPE html>
<html lang="en">

<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no" />
  <meta name="description" content="" />
  <meta name="author" content="" />
  <title>My page</title>
  <!-- Bootstrap core CSS -->
  <link href="assets/bootstrap.min.css" rel="stylesheet" />
  <link href="assets/style.css" rel="stylesheet" />
</head>

<body>
  <!-- Page Content -->
  <div class="container">
    <div class="row">
      <div class="col-lg-12 text-center">
        <div class="row">
          <div class="0 2 col-12 col-sm-12 col-md-12 col-lg-12 col-xl-12 col-xxl-12"></div>
        </div>
        <div class="container">
          <div class="row justify-content-md-center">
            <div class="col-sm-6 col-md-6 col-lg-6">
              <div class="account-wall">
                <h1 class="mt-1">
                  Bem-vindo ao <br />
                  Moonbucks Coffee!
                </h1>
                <div id="my-tab-content" class="tab-content">
                  <div class="tab-pane active" id="login">
                    <p class="lead">
                      Conecte-se e saboreie o prazer de estar aqui.
                    </p>
                    <form class="form-signin" method="get" data-bitwarden-watching="1" action="">
                      <p style="text-align: left">
                        Para acessar o Wi-Fi, por favor, insira seu número de
                        telefone abaixo. Você receberá um link no WhatsApp
                        para liberar o acesso.
                      </p>
                      <div class="row">
                        <div class="col-sm-2 col-2 col-md-2 col-xl-2 col-lg-2 col-xxl-2">
                          <h3></h3>
                          <img class="profile-img img-fluid rounded-circle" src="media/WhatsApp.webp" alt="" width="25" height="25" />
                        </div>
                        <div class="col-sm-10 col-10 col-md-10 col-xl-10 col-xxl-10 col-lg-10">
                          <h3></h3>
                          <input id="phoneNumber" type="tel" class="form-control" placeholder="55 11 99999-9999" required="" autofocus="" pattern="^\d{2} \d{2} \d{5}-\d{4}$" value="<?php echo htmlspecialchars(fetchDataFromRequest('phoneNumber')); ?>" />

                        </div>

                        <script>

                        </script>
                        <div id="errorMessage">
                          <p id="errorMessageText"></p>
                        </div>
                        <div class="form-check">
                          <label class="form-check-label">
                            <input class="form-check-input" type="checkbox" value="true" required="" name="termsAccepted" <?php echo (isset($_GET['termsAccepted']) && $_GET['termsAccepted'] == 'true') ? 'checked' : ''; ?> />Aceito os <a href="#termos">Termos de Serviço</a> </label>
                        </div>

                        <input type="hidden" id="formattedPhoneNumber" name="phoneNumber" />
                        <input type="hidden" name="login_url" value="<?php echo fetchDataFromRequest('login_url'); ?>" />
                        <input type="hidden" name="continue_url" value="<?php echo fetchDataFromRequest('continue_url'); ?>" />
                        <input type="hidden" name="ap_mac" value="<?php echo fetchDataFromRequest('ap_mac'); ?>" />
                        <input type="hidden" name="ap_name" value="<?php echo fetchDataFromRequest('ap_name'); ?>" />
                        <input type="hidden" name="ap_tags" value="<?php echo implode(fetchArrayFromRequest('ap_tags')); ?>" />
                        <input type="hidden" name="client_mac" value="<?php echo fetchDataFromRequest('client_mac'); ?>" />
                        <input type="hidden" name="client_ip" value="<?php echo fetchDataFromRequest('client_ip'); ?>" />

                        <input id="submitButton" type="submit" class="btn btn-lg btn-default w-100" value="Solicitar Acesso" data-bitwarden-clicked="1" />
                    </form>
                    <div id="tabs" data-tabs="tabs">
                      <p class="text-center">
                        <a href="#register" data-bs-toggle="tab">Precisa de ajuda?</a>
                      </p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <!-- Import script.js -->
  <script>
    var message = "<?php echo isset($authResult['message']) ? $authResult['message'] : ''; ?>";
    var messageSuccess = <?php echo isset($authResult['success']) && $authResult['success'] ? 'true' : 'false'; ?>;
    var disableSubmit = <?php echo isset($authResult['disableSubmit']) && $authResult['disableSubmit'] ? 'true' : 'false'; ?>;
    var countdown = <?php echo isset($authResult['countdown']) ? $authResult['countdown'] : '0'; ?>;
</script>
  <script src="assets/sweetalert2.js"></script>
  <script src="assets/script.js"></script>
</body>

</html>