
import * as React from 'react';
import CssBaseline from '@mui/material/CssBaseline';
import Box from '@mui/material/Box';
import Container from '@mui/material/Container';
import Grid from '@mui/material/Grid';
import Card from '@mui/material/Card';
import CardActions from '@mui/material/CardActions';
import CardContent from '@mui/material/CardContent';
import CardMedia from '@mui/material/CardMedia';
import Button from '@mui/material/Button';
import Typography from '@mui/material/Typography';

export default function Home({ data }) {
  return (
    <>
      <div>
      <Container maxWidth="lg">
        <Grid container spacing={2}>

          {data.map((item) => (
            <Grid item xs={12} md={4}>
              <Card>
                <CardMedia
                  component="img"
                  height="140"
                  image={item.coverimage}
                  alt="green iguana"
                />
                <CardContent>
                  <Typography gutterBottom variant="h5" component="div">
                  {item.name}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                  {item.detail}
                  </Typography>
                </CardContent>
                <CardActions>
                  <Button size="small">Share</Button>
                  <Button size="small">Learn More</Button>
                </CardActions>
              </Card>
            </Grid>
          ))}

        </Grid>
      </Container>
      </div>
    </>
  )
}

export async function getServerSideProps() {
  const res = await fetch('http://localhost:8080/attractions')
  const data = await res.json()
  console.log(data)

  return { props: { data } }
}

