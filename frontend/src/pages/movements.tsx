import { useEffect, useState, useCallback } from 'react';
import { api, ApiError } from '../api/client';
import type { StockMovement } from '../types/index';
import { formatRelativeTime, formatDate } from '../utils/format';
import '../App.css';

const StockMovements = () => {
  const [movements, setMovements] = useState<StockMovement[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [search, setSearch] = useState('');
  const [typeFilter, setTypeFilter] = useState('');
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [pageSize] = useState(20);

  const fetchMovements = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await api.getMovements(search, typeFilter, page, pageSize);
      setMovements(data.data);
      setTotal(data.total);
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      } else {
        setError('Failed to load stock movements');
      }
    } finally {
      setLoading(false);
    }
  }, [search, typeFilter, page, pageSize]);

  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    fetchMovements();
  }, [fetchMovements]);

  const totalPages = Math.ceil(total / pageSize);

  return (
    <>
      <div className="page-header">
        <h2>Stock Movements</h2>
        <p>View all stock adjustments and transactions</p>
      </div>

      <div className="toolbar" style={{ marginBottom: '20px' }}>
        <input
          type="text"
          className="search-input"
          placeholder="Search by SKU or product name..."
          value={search}
          onChange={(e) => {
            setSearch(e.target.value);
            setPage(1);
          }}
        />
        <select
          className="form-control"
          value={typeFilter}
          onChange={(e) => {
            setTypeFilter(e.target.value);
            setPage(1);
          }}
          style={{ flex: '0 1 120px' }}
        >
          <option value="">All Types</option>
          <option value="in">Stock In</option>
          <option value="out">Stock Out</option>
        </select>
      </div>

      {error && <div className="error-message">{error}</div>}

      {loading ? (
        <div className="loading">Loading movements...</div>
      ) : movements.length === 0 ? (
        <div className="empty-state">
          <p>No stock movements</p>
          <small>Stock adjustments will appear here</small>
        </div>
      ) : (
        <div className="card">
          <div className="table-container">
            <table>
              <thead>
                <tr>
                  <th>Date</th>
                  <th>SKU</th>
                  <th>Product Name</th>
                  <th>Type</th>
                  <th>Qty</th>
                  <th>Note</th>
                </tr>
              </thead>
              <tbody>
                {movements.map((movement) => (
                  <tr key={movement.id}>
                    <td>
                      <div style={{ fontSize: '13px' }}>
                        <div>{formatDate(movement.created_at)}</div>
                        <div style={{ color: '#64748b' }}>{formatRelativeTime(movement.created_at)}</div>
                      </div>
                    </td>
                    <td>{movement.sku}</td>
                    <td>{movement.product_name}</td>
                    <td>
                      <span className={`badge ${movement.type === 'in' ? 'badge-success' : 'badge-danger'}`}>
                        {movement.type.toUpperCase()}
                      </span>
                    </td>
                    <td>{movement.qty}</td>
                    <td style={{ fontSize: '13px', color: '#64748b' }}>{movement.note || '-'}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {totalPages > 1 && (
            <div className="pagination">
              <button onClick={() => setPage(1)} disabled={page === 1}>
                First
              </button>
              <button onClick={() => setPage(page - 1)} disabled={page === 1}>
                Previous
              </button>
              <span style={{ padding: '8px 12px', color: '#1e293b' }}>
                Page {page} of {totalPages}
              </span>
              <button onClick={() => setPage(page + 1)} disabled={page === totalPages}>
                Next
              </button>
              <button onClick={() => setPage(totalPages)} disabled={page === totalPages}>
                Last
              </button>
            </div>
          )}
        </div>
      )}
    </>
  );
};

export default StockMovements;
