import { LocalShipping } from "@mui/icons-material";
import { AppBar, IconButton, Toolbar, Typography } from "@mui/material";

export const Navbar = () => {
    return (
        <AppBar position="static">
            <Toolbar>
                <IconButton onClick={() => window.location.href = "/"} size="large" edge='start' color="inherit" aria-label="logo">
                    <LocalShipping />
                </IconButton>
                <Typography variant="h6" component='div'>
                    Starter Project
                </Typography>
            </Toolbar>
        </AppBar>
    )
}