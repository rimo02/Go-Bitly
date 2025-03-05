import { useState, useEffect } from 'react';
import axios from 'axios';
import { TableContainer, Table, TableHead, TableBody, TableRow, TableCell, Paper, Typography, Container } from '@mui/material';

function Home() {
  const [formData, setFormData] = useState([]);
  const [error, setError] = useState("");

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get('http://localhost:5000/dashboard');
        setFormData(response.data);
      } catch (err) {
        setError(err.message);
      }
    };
    fetchData();
  }, []);

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
        bgcolor: "#121212",
        color: "#fff",
        padding: 4,
      }}
    >
      <Typography variant="h4" sx={{ marginBottom: 3, fontWeight: "bold", textAlign: "center" }}>
        🔗 Shortened URLs Dashboard
      </Typography>

      {error && <Typography color="error" sx={{ marginBottom: 2 }}>{error}</Typography>}

      <TableContainer
        component={Paper}
        sx={{
          maxWidth: 800,
          bgcolor: "#1E1E1E",
          borderRadius: 2,
          boxShadow: "0px 4px 12px rgba(0,0,0,0.3)",
        }}
      >
        <Table sx={{ minWidth: 600}}>
          <TableHead>
            <TableRow sx={{ backgroundColor: "#333" }}>
              {formData.length > 0 &&
                Object.keys(formData[0]).map((key, index) =>
                  index !== 0 ? (
                    <TableCell key={index} sx={{ color: "#fff", fontWeight: "bold" }}>
                      {key.toUpperCase()}
                    </TableCell>
                  ) : null
                )}
            </TableRow>
          </TableHead>

          <TableBody>
            {formData.map((row, index) => (
              <TableRow
                key={index}
                sx={{
                  "&:nth-of-type(odd)": { backgroundColor: "#252525" },
                  "&:hover": { backgroundColor: "#444" }
                }}
              >
                {Object.values(row).map((value, i) => {
                  let displayValue;

                  if (i === 4 || i === 5) {
                    displayValue = new Date(value * 1000).toLocaleDateString();
                  } else if (i === 2) {
                    displayValue = 'localhost:5000/' + value
                  } else {
                    displayValue = value;
                  }

                  return (
                    i !== 0 && (
                      <TableCell key={i} sx={{ color: "#ddd", padding: "12px" }}>
                        {displayValue}
                      </TableCell>
                    )
                  );
                })}

              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </Container>
  );
}

export default Home;
