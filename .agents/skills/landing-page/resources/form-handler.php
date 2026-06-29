<?php
/**
 * Simple Secure Form Handler for Landing Pages
 * Usage: Set your form action to this file (e.g., action="form-handler.php") and method to POST.
 */

session_start();

// Configuration
$recipient_email = "your-email@example.com";
$subject_prefix = "[Landing Page Lead]";
$success_redirect = "index.php?status=success";
$error_redirect = "index.php?status=error";

// Basic CSRF Protection (Generate token in your HTML form if possible)
// <input type="hidden" name="csrf_token" value="<?php echo $_SESSION['csrf_token']; ?>">

if ($_SERVER["REQUEST_METHOD"] == "POST") {
    
    // Sanitize Inputs
    $name = filter_input(INPUT_POST, 'name', FILTER_SANITIZE_STRING);
    $email = filter_input(INPUT_POST, 'email', FILTER_SANITIZE_EMAIL);
    $message = filter_input(INPUT_POST, 'message', FILTER_SANITIZE_STRING);
    
    // Validate required fields
    if (empty($name) || empty($email) || !filter_var($email, FILTER_VALIDATE_EMAIL)) {
        header("Location: $error_redirect&reason=validation");
        exit;
    }
    
    // Prepare Email
    $subject = "$subject_prefix New submission from $name";
    $body = "Name: $name\n";
    $body .= "Email: $email\n\n";
    $body .= "Message:\n$message\n";
    
    $headers = "From: no-reply@yourdomain.com\r\n";
    $headers .= "Reply-To: $email\r\n";
    
    // Send Email
    if (mail($recipient_email, $subject, $body, $headers)) {
        header("Location: $success_redirect");
    } else {
        header("Location: $error_redirect&reason=server");
    }
    exit;
} else {
    // Not a POST request
    header("Location: index.php");
    exit;
}
?>
