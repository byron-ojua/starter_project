import { Container } from "@mui/material";
import { useEffect } from "react";

/**
 * Page that displays a 404 error message.
 * @returns [JSX.Element] 404 page
 */
const Error404 = () => {
    useEffect(() => {
        document.title = "404 | Starter Project"
    }, [])

    return (
        <div className="App">
            <Container>
                <h1>404</h1>
                <h2>Page not found</h2>
                <p>The page you are looking for does not exist.</p>
            </Container>
        </div>
    )
};

export default Error404;