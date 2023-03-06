<?php
// index.php

// 如果 URL 请求带参数 health，则返回状态码 200
if (isset($_GET['health'])) {
    http_response_code(200);
    echo "OK\n";
    exit;
}

echo gethostname();
date_default_timezone_set("Asia/Shanghai");
echo date("l jS \of F Y h:i:s A");
sleep(1);
echo " v1\n";
?>
