import './App.css';
import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { About, Home } from './components/pages'
import { Navbar } from './components/navbar'
import { Mails } from './components/mails'
import { Auth } from './components/auth'
import { NewMail } from './components/new-mail';
import { Currencies } from './components/currensies';


function App() {
  return (
    <main>
    <Navbar />
    <BrowserRouter>
      <div>
        <Routes>
          <Route path='/' element={<Home />} />
          <Route path='/about' element={<About />} />
          <Route path='/mails' element={<Mails />} />
          <Route path='/login' element={<Auth />} />
          <Route path='/mails/new' element={<NewMail />} />
          <Route path='/currencies' element={<Currencies />} />
        </Routes>
      </div>
    </BrowserRouter>
  </main>
  );
}

export default App;
