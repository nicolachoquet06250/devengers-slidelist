export default function () {
    const { token } = window.params;

    if (token !== "") {
        localStorage.setItem("googleOAuthToken", token);
        window.close();
    }
}