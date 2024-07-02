<?php

function fetchDataFromRequest($key)
{
    return isset($_REQUEST[$key]) ? urldecode($_REQUEST[$key]) : '';
}

function fetchArrayFromRequest($key)
{
    return isset($_REQUEST[$key]) ? explode(" ", $_REQUEST[$key]) : [];
}
function generateUserData()
{
    return [
        'phoneNumber' => fetchDataFromRequest('phoneNumber'),
        'loginURL' => fetchDataFromRequest('login_url'),
        'continueURL' => fetchDataFromRequest('continue_url'),
        'apMAC' => fetchDataFromRequest('ap_mac'),
        'apName' => fetchDataFromRequest('ap_name'),
        'apTags' => fetchArrayFromRequest('ap_tags'),
        'clientMAC' => fetchDataFromRequest('client_mac'),
        'clientIP' => fetchDataFromRequest('client_ip'),
    ];
}


function sendCurlRequest($url, $requestData)
{
    $ch = curl_init($url);
    curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, false);
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
    curl_setopt($ch, CURLOPT_POSTFIELDS, json_encode($requestData));
    curl_setopt($ch, CURLOPT_HTTPHEADER, ['Content-Type: application/json']);

    $response = curl_exec($ch);
    if (!$response) {
        throw new Exception('CURL Error: ' . curl_error($ch), curl_errno($ch));
    }

    curl_close($ch);
    return json_decode($response, true);
}

function processAuthentication()
{
    $accountManagerUrl = getenv('ACCOUNTMANAGER') ?: 'http://ACCOUNTMANAGER';
    $requestData = generateUserData();
    $responseData = sendCurlRequest($accountManagerUrl . '/api/user', $requestData);

    if ($responseData['success']) {
        return ['message' => $responseData['message'], 'success' => true, 'disableSubmit' => true, 'countdown' => 120];
    } else {
        return ['message' => $responseData['message'], 'success' => false, 'disableSubmit' => false, 'countdown' => 0];
    }
}
?>