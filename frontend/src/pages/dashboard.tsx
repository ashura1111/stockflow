import { useEffect, useState, useCallback } from 'react';
import { api, ApiError } from '../api/client';
import type { DashboardStats } from '../types/index';
import { formatCurrency, formatRelativeTime } from '../utils/format';
import '../App.css';

const Dashboard = () => {
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchStats = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await api.dashboard();
      setStats(data);
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      } else {
        setError('Failed to load dashboard');
      }
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    fetchStats();
  }, [fetchStats]);

  if (loading) return <div className="loading">Loading dashboard...</div>;
  if (error) return <div className="error-message">{error}</div>;
  if (!stats) return null;

  return (
    <>
      <div className="page-header">
        <h2>Dashboard</h2>
        <p>Overview of your inventory and sales</p>
      </div>

      <div className="stats-grid">
        <div className="stat-card">
          <h3>Total Products</h3>
          <div className="value">{stats.total_products}</div>
        </div>
        <div className="stat-card">
          <h3>Active Products</h3>
          <div className="value">{stats.active_products}</div>
        </div>
        <div className="stat-card">
          <h3>Total Stock Value</h3>
          <div className="value">{formatCurrency(stats.total_stock_value)}</div>
        </div>
        <div className="stat-card">
          <h3>Low Stock Products</h3>
          <div className="value">{stats.low_stock_products}</div>
        </div>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(400px, 1fr))', gap: '20px' }}>
        <div className="card">
          <h3>Recent Stock Movements</h3>
          {stats.recent_movements.length === 0 ? (
            <p style={{ color: '#64748b', marginTop: '12px' }}>No recent movements</p>
          ) : (
            <table style={{ marginTop: '12px', fontSize: '13px' }}>
              <tbody>
                {stats.recent_movements.slice(0, 5).map((movement) => (
                  <tr key={movement.id}>
                    <td style={{ padding: '8px' }}>
                      <span className={`badge ${movement.type === 'in' ? 'badge-success' : 'badge-danger'}`}>
                        {movement.type.toUpperCase()}
                      </span>
                    </td>
                    <td style={{ padding: '8px' }}>{movement.sku}</td>
                    <td style={{ padding: '8px' }}>{movement.qty}</td>
                    <td style={{ padding: '8px', color: '#64748b' }}>
                      {formatRelativeTime(movement.created_at)}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>

        <div className="card">
          <h3>Recent Orders</h3>
          {stats.recent_orders.length === 0 ? (
            <p style={{ color: '#64748b', marginTop: '12px' }}>No recent orders</p>
          ) : (
            <table style={{ marginTop: '12px', fontSize: '13px' }}>
              <tbody>
                {stats.recent_orders.slice(0, 5).map((order) => (
                  <tr key={order.id}>
                    <td style={{ padding: '8px' }}>Order #{order.id}</td>
                    <td style={{ padding: '8px' }}>{formatCurrency(order.total)}</td>
                    <td style={{ padding: '8px' }}>
                      <span className={`badge ${order.status === 'CREATED' ? 'badge-info' : 'badge-danger'}`}>
                        {order.status}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      </div>

      {stats.low_stock_list.length > 0 && (
        <div className="card" style={{ marginTop: '20px' }}>
          <h3>Low Stock Products</h3>
          <table style={{ marginTop: '12px', fontSize: '13px' }}>
            <thead>
              <tr>
                <th>SKU</th>
                <th>Name</th>
                <th>Stock</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              {stats.low_stock_list.map((product) => (
                <tr key={product.id}>
                  <td>{product.sku}</td>
                  <td>{product.name}</td>
                  <td>{product.stock}</td>
                  <td>
                    <span className="badge badge-warning">Low Stock</span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </>
  );
};

export default Dashboard;
