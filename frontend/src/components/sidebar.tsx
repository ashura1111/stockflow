import { Link, useLocation } from 'react-router-dom';
import { 
  LayoutDashboard, 
  Package, 
  Activity, 
  ShoppingCart
} from 'lucide-react';
import '../App.css';

const Sidebar = () => {
  const location = useLocation();

  const menuItems = [
    { path: '/', icon: LayoutDashboard, label: 'Dashboard' },
    { path: '/products', icon: Package, label: 'Products' },
    { path: '/movements', icon: Activity, label: 'Stock Movements' },
    { path: '/orders', icon: ShoppingCart, label: 'Sales Orders' },
  ];

  return (
    <nav className="sidebar">
      <div className="sidebar-header">
        <h1>StockFlow</h1>
      </div>
      <div className="sidebar-nav">
        {menuItems.map((item) => {
          const Icon = item.icon;
          const isActive = location.pathname === item.path;
          return (
            <Link key={item.path} to={item.path} className={`nav-item ${isActive ? 'active' : ''}`}>
              <Icon size={20} />
              {item.label}
            </Link>
          );
        })}
      </div>
    </nav>
  );
};

export default Sidebar;
