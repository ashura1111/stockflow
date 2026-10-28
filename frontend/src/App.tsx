import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Layout from './components/layout';
import Dashboard from './pages/dashboard';
import Products from './pages/products';
import StockMovements from './pages/movements';
import SalesOrders from './pages/orders';
import './App.css';

function App() {
  return (
    <Router>
      <Layout>
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/products" element={<Products />} />
          <Route path="/movements" element={<StockMovements />} />
          <Route path="/orders" element={<SalesOrders />} />
        </Routes>
      </Layout>
    </Router>
  );
}

export default App;
