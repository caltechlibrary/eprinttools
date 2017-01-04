<?php

function httpGET($url) {
    if (function_exists("curl_init")) {
        $ch = curl_init(); 
        curl_setopt($ch, CURLOPT_URL, $url);
        curl_setopt($ch, CURLOPT_HEADER, 0);
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, TRUE);
        // Set timeput for 30 seconds.
        curl_setopt($ch, CURLOPT_TIMEOUT, 30);
        if ( ! $result = curl_exec($ch)) { 
            trigger_error(curl_error($ch)); 
        } 
        curl_close($ch); 
        return $result; 
    }
    return file_get_content($url);
}

// Get the articles from the website, if open urls not allowed use the CURL PHP Package to fetch the content
function getSomeArticles($url, $maxLength = 5) {
    $data = httpGET($url);
    $articles = json_decode($data);
    $results = array();
    for ($i = 0; $i < count($articles) && $i < $maxLength; $i++) {
	    array_push($results, $articles[$i]);
    }
    return $results;
}

$baseURL = getenv("EPGO_SITE_URL");
if ($baseURL === "") {
    if (php_sapi_name() !== "cli") {
        header("HTTP/1.0 404 Not Found");
    }
    exit(1);
} else {
    $myShortList = getSomeArticles("$baseURL/recent/articles.json", 3);
    if (php_sapi_name() !== "cli") {
        header("Content-Type: application/json");
    }
}
echo json_encode($myShortList, JSON_PRETTY_PRINT);
?>

