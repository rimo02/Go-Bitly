import { useState } from 'react';
import { Container, Button, Typography, Box, TextField, Paper, IconButton, AppBar, Toolbar } from '@mui/material';
import { useNavigate } from 'react-router-dom'
import axios from 'axios';
import ContentCopyIcon from '@mui/icons-material/ContentCopy';

function ShortenPage() {
    const [originalUrl, setOriginalUrl] = useState("");
    const [shortUrl, setShortUrl] = useState("");
    const [error, setError] = useState("");
    const navigate = useNavigate()
    const handleShorten = async () => {
        try {
            const response = await axios.post('http://localhost:5000/shorten', {
                lurl: originalUrl,
                hrs: 0,
                mins: 0,
                days: 2,
            });
            setShortUrl(response.data.ShortUrl);
            setError("");
        } catch (err) {
            setError("Failed to shorten URL");
            console.log(err)
            setShortUrl("");
        }
    };

    const handleCopy = () => {
        navigator.clipboard.writeText(shortUrl);
    };

    return (
        <Container
            maxWidth={false}
            disableGutters
            sx={{
                display: "flex",
                flexDirection: "column",
                alignItems: "center",
                justifyContent: "center",
                height: "100vh",
                width: "100vw",
                bgcolor: "#191919FF",
                color: "#fff",
            }}
        >
            <AppBar position="absolute" sx={{ bgcolor: "#1E1E1EFF" }}>
                <Toolbar sx={{ display: "flex", justifyContent: "flex-end" }}>
                    <Button color="inherit" onClick={() => navigate('/home')}>Home</Button>
                </Toolbar>
            </AppBar>

            <Paper
                elevation={20}
                sx={{
                    padding: 4,
                    width: "80%",
                    maxWidth: "500px",
                    bgcolor: "#373737FF",
                    textAlign: "center",
                    borderRadius: 3
                }}
            >
                <Typography variant="h4" gutterBottom sx={{ fontWeight: "bold", color: "#90caf9" }}>
                    URL Shortener
                </Typography>
                <TextField
                    label="Enter URL"
                    variant="outlined"
                    fullWidth
                    value={originalUrl}
                    onChange={(e) => setOriginalUrl(e.target.value)}
                    sx={{ marginBottom: 2, input: { color: "white" }, label: { color: "gray" } }}
                />
                <Button
                    variant="contained"
                    color="primary"
                    fullWidth
                    onClick={handleShorten}
                    sx={{ marginBottom: 2, bgcolor: "#1976d2", ":hover": { bgcolor: "#1565c0" }, width: "50%", maxWidth: "200px" }}
                >
                    Shorten
                </Button>

                {shortUrl && (
                    <Box sx={{ display: "flex", alignItems: "center", justifyContent: "center", bgcolor: "#333", p: 1, borderRadius: 2 }}>
                        <Typography variant="body1" sx={{ color: "#90caf9", marginRight: 1 }}>
                            <a href={shortUrl} target="_blank" rel="noopener noreferrer" style={{ color: "#90caf9", textDecoration: "none" }}>
                                {shortUrl}
                            </a>
                        </Typography>
                        <IconButton onClick={handleCopy} sx={{ color: "white" }}>
                            <ContentCopyIcon />
                        </IconButton>
                    </Box>
                )}
                {error && <Typography color="error">{error}</Typography>}
            </Paper>
        </Container>
    );
}

export default ShortenPage