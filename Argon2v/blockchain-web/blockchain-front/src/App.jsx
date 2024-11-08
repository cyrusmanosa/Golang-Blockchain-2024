import Login from './login'
import Check from './check'
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Login />} />
        <Route path="/Ok/:name" element={<Check />} />
      </Routes>
    </Router>
  )
}

export default App
