import { useEffect, useState, useCallback } from 'react';
import { Eye, Trash2, Plus } from 'lucide-react';
import { api, ApiError } from '../api/client';
import type { SalesOrder, CreateOrderRequest, Product } from '../types/index';
import { formatCurrency, formatDate } from '../utils/format';
import Modal from '../components/modal';
import Toast from '../components/toast';
import '../App.css';

interface OrderLineItem {
  sku: string;
  qty: number;
}

const SalesOrders = () => {
  const [orders, setOrders] = useState<SalesOrder[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [pageSize] = useState(10);

  const [showCreateModal, setShowCreateModal] = useState(false);
  const [customerName, setCustomerName] = useState('');
  const [orderItems, setOrderItems] = useState<OrderLineItem[]>([]);
  const [availableProducts, setAvailableProducts] = useState<Product[]>([]);
  const [createError, setCreateError] = useState<string | null>(null);
  const [createSubmitting, setCreateSubmitting] = useState(false);

  const [showOrderDetail, setShowOrderDetail] = useState(false);
  const [selectedOrder, setSelectedOrder] = useState<SalesOrder | null>(null);

  const [showConfirmCancel, setShowConfirmCancel] = useState(false);
  const [cancellingOrder, setCancellingOrder] = useState<SalesOrder | null>(null);
  const [cancelSubmitting, setCancelSubmitting] = useState(false);

  const [toast, setToast] = useState<{ message: string; type: 'success' | 'error' } | null>(null);

  const fetchOrders = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await api.getOrders('', '', page, pageSize);
      setOrders(data.data);
      setTotal(data.total);
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      } else {
        setError('Failed to load orders');
      }
    } finally {
      setLoading(false);
    }
  }, [page, pageSize]);

  const fetchAvailableProducts = useCallback(async () => {
    try {
      const data = await api.getProducts('', 'active', false, 1, 1000);
      setAvailableProducts(data.data);
    } catch {
      setAvailableProducts([]);
    }
  }, []);

  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    fetchOrders();
  }, [fetchOrders]);

  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    fetchAvailableProducts();
  }, [fetchAvailableProducts]);

  const handleOpenCreateModal = () => {
    setCustomerName('');
    setOrderItems([{ sku: '', qty: 1 }]);
    setCreateError(null);
    setShowCreateModal(true);
  };

  const handleAddOrderItem = () => {
    setOrderItems([...orderItems, { sku: '', qty: 1 }]);
  };

  const handleRemoveOrderItem = (index: number) => {
    setOrderItems(orderItems.filter((_, i) => i !== index));
  };

  const handleOrderItemChange = (index: number, field: string, value: string | number) => {
    const updated = [...orderItems];
    updated[index] = { ...updated[index], [field]: value };
    setOrderItems(updated);
  };

  const handleCreateOrder = async () => {
    setCreateError(null);

    if (!customerName.trim()) {
      setCreateError('Customer name is required');
      return;
    }

    if (orderItems.length === 0) {
      setCreateError('At least one item is required');
      return;
    }

    const skus = orderItems.map((item) => item.sku);
    if (new Set(skus).size !== skus.length) {
      setCreateError('Duplicate SKU in order');
      return;
    }

    const request: CreateOrderRequest = {
      customer: customerName,
      items: orderItems.map((item) => ({
        sku: item.sku,
        qty: item.qty,
      })),
    };

    setCreateSubmitting(true);

    try {
      const createdOrder = await api.createOrder(request);
      setToast({ message: 'Order created successfully', type: 'success' });
      setShowCreateModal(false);
      setSelectedOrder(createdOrder);
      setShowOrderDetail(true);
      await fetchOrders();
    } catch (err) {
      if (err instanceof ApiError) {
        if (err.code === 'INSUFFICIENT_STOCK') {
          setCreateError(`Insufficient stock for ${err.details.sku}. Available: ${err.details.available}`);
        } else {
          setCreateError(err.message);
        }
      } else {
        setCreateError('Failed to create order');
      }
    } finally {
      setCreateSubmitting(false);
    }
  };

  const handleOpenOrderDetail = async (order: SalesOrder) => {
    try {
      const fullOrder = await api.getOrder(order.id);
      setSelectedOrder(fullOrder);
      setShowOrderDetail(true);
    } catch {
      setToast({ message: 'Failed to load order details', type: 'error' });
    }
  };

  const handleOpenCancelConfirm = (order: SalesOrder) => {
    setCancellingOrder(order);
    setShowConfirmCancel(true);
  };

  const handleConfirmCancel = async () => {
    if (!cancellingOrder) return;

    setCancelSubmitting(true);

    try {
      await api.cancelOrder(cancellingOrder.id);
      setToast({ message: 'Order cancelled successfully', type: 'success' });
      setShowConfirmCancel(false);
      setShowOrderDetail(false);
      await fetchOrders();
    } catch (err) {
      if (err instanceof ApiError) {
        setToast({ message: err.message, type: 'error' });
      } else {
        setToast({ message: 'Failed to cancel order', type: 'error' });
      }
      setShowConfirmCancel(false);
    } finally {
      setCancelSubmitting(false);
    }
  };

  const totalPages = Math.ceil(total / pageSize);
  const calculateOrderTotal = (items: OrderLineItem[]): number => {
    return items.reduce((total, item) => {
      const product = availableProducts.find((p) => p.sku === item.sku);
      return total + (product ? product.price * item.qty : 0);
    }, 0);
  };

  return (
    <>
      <div className="page-header">
        <h2>Sales Orders</h2>
        <p>Manage customer orders</p>
      </div>

      {error && <div className="error-message">{error}</div>}

      <div style={{ marginBottom: '20px' }}>
        <button className="btn btn-primary" onClick={handleOpenCreateModal}>
          <Plus size={16} /> Create Order
        </button>
      </div>

      {loading ? (
        <div className="loading">Loading orders...</div>
      ) : orders.length === 0 ? (
        <div className="empty-state">
          <p>No sales orders</p>
          <small>Create your first order to get started</small>
        </div>
      ) : (
        <div className="card">
          <div className="table-container">
            <table>
              <thead>
                <tr>
                  <th>Order ID</th>
                  <th>Customer</th>
                  <th>Date</th>
                  <th>Total</th>
                  <th>Status</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {orders.map((order) => (
                  <tr key={order.id}>
                    <td>#{order.id}</td>
                    <td>{order.customer}</td>
                    <td>{formatDate(order.created_at)}</td>
                    <td>{formatCurrency(order.total)}</td>
                    <td>
                      <span className={`badge ${order.status === 'CREATED' ? 'badge-info' : 'badge-danger'}`}>
                        {order.status}
                      </span>
                    </td>
                    <td>
                      <div style={{ display: 'flex', gap: '8px' }}>
                        <button
                          className="btn btn-secondary"
                          onClick={() => handleOpenOrderDetail(order)}
                        >
                          <Eye size={14} />
                        </button>
                        {order.status === 'CREATED' && (
                          <button
                            className="btn btn-danger"
                            onClick={() => handleOpenCancelConfirm(order)}
                          >
                            <Trash2 size={14} />
                          </button>
                        )}
                      </div>
                    </td>
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

      <Modal
        isOpen={showCreateModal}
        onClose={() => setShowCreateModal(false)}
        title="Create Sales Order"
        onSubmit={handleCreateOrder}
        submitLabel="Create Order"
        isLoading={createSubmitting}
      >
        {createError && <div className="error-message" style={{ marginBottom: '12px' }}>{createError}</div>}

        <div className="form-group">
          <label>Customer Name</label>
          <input
            type="text"
            className="form-control"
            value={customerName}
            onChange={(e) => setCustomerName(e.target.value)}
            placeholder="e.g., PT Contoh Properti"
          />
        </div>

        <div style={{ marginTop: '20px', marginBottom: '12px' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '12px' }}>
            <label style={{ margin: 0 }}>Items</label>
            <button className="btn btn-secondary" onClick={handleAddOrderItem} type="button">
              <Plus size={14} /> Add Item
            </button>
          </div>

          {orderItems.map((item, index) => (
            <div key={index} style={{ display: 'flex', gap: '8px', marginBottom: '8px' }}>
              <select
                className="form-control"
                value={item.sku}
                onChange={(e) => handleOrderItemChange(index, 'sku', e.target.value)}
                style={{ flex: 1 }}
              >
                <option value="">Select Product</option>
                {availableProducts.map((product) => (
                  <option key={product.id} value={product.sku}>
                    {product.sku} - {product.name} ({product.stock} available)
                  </option>
                ))}
              </select>
              <input
                type="number"
                className="form-control"
                value={item.qty}
                onChange={(e) => handleOrderItemChange(index, 'qty', parseInt(e.target.value) || 0)}
                placeholder="Qty"
                min="1"
                style={{ width: '80px' }}
              />
              <button
                className="btn btn-danger"
                onClick={() => handleRemoveOrderItem(index)}
                type="button"
                style={{ minWidth: 'auto' }}
              >
                <Trash2 size={14} />
              </button>
            </div>
          ))}
        </div>

        {orderItems.length > 0 && (
          <div style={{ 
            backgroundColor: '#f8fafc', 
            padding: '12px', 
            borderRadius: '6px',
            marginTop: '12px'
          }}>
            <div style={{ display: 'flex', justifyContent: 'space-between' }}>
              <strong>Estimated Total:</strong>
              <strong>{formatCurrency(calculateOrderTotal(orderItems))}</strong>
            </div>
          </div>
        )}
      </Modal>

      {showOrderDetail && selectedOrder && (
        <div className="modal-overlay" onClick={() => setShowOrderDetail(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h3>Order #{selectedOrder.id} - {selectedOrder.customer}</h3>
            </div>
            <div className="modal-body">
              <div style={{ marginBottom: '16px' }}>
                <p><strong>Created:</strong> {formatDate(selectedOrder.created_at)}</p>
                <p><strong>Status:</strong> <span className={`badge ${selectedOrder.status === 'CREATED' ? 'badge-info' : 'badge-danger'}`}>{selectedOrder.status}</span></p>
                {selectedOrder.cancelled_at && <p><strong>Cancelled:</strong> {formatDate(selectedOrder.cancelled_at)}</p>}
              </div>

              <table style={{ marginBottom: '16px' }}>
                <thead>
                  <tr>
                    <th>SKU</th>
                    <th>Product</th>
                    <th>Qty</th>
                    <th>Unit Price</th>
                    <th>Subtotal</th>
                  </tr>
                </thead>
                <tbody>
                  {selectedOrder.items.map((item) => (
                    <tr key={item.id}>
                      <td>{item.sku}</td>
                      <td>{item.product_name}</td>
                      <td>{item.qty}</td>
                      <td>{formatCurrency(item.unit_price)}</td>
                      <td>{formatCurrency(item.subtotal)}</td>
                    </tr>
                  ))}
                </tbody>
              </table>

              <div style={{ backgroundColor: '#f8fafc', padding: '12px', borderRadius: '6px', marginBottom: '16px' }}>
                <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                  <strong>Order Total:</strong>
                  <strong>{formatCurrency(selectedOrder.total)}</strong>
                </div>
              </div>
            </div>
            <div className="modal-footer">
              <button className="btn btn-secondary" onClick={() => setShowOrderDetail(false)}>
                Close
              </button>
              {selectedOrder.status === 'CREATED' && (
                <button className="btn btn-danger" onClick={() => handleOpenCancelConfirm(selectedOrder)}>
                  <Trash2 size={14} /> Cancel Order
                </button>
              )}
            </div>
          </div>
        </div>
      )}

      {showConfirmCancel && cancellingOrder && (
        <div className="modal-overlay" onClick={() => setShowConfirmCancel(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h3>Confirm Cancellation</h3>
            </div>
            <div className="modal-body">
              <p>Are you sure you want to cancel order <strong>#{cancellingOrder.id}</strong>?</p>
              <p style={{ fontSize: '13px', color: '#64748b', marginTop: '12px' }}>
                Stock will be restored for all items. This action cannot be undone.
              </p>
            </div>
            <div className="modal-footer">
              <button className="btn btn-secondary" onClick={() => setShowConfirmCancel(false)} disabled={cancelSubmitting}>
                Keep Order
              </button>
              <button className="btn btn-danger" onClick={handleConfirmCancel} disabled={cancelSubmitting}>
                {cancelSubmitting ? 'Cancelling...' : 'Cancel Order'}
              </button>
            </div>
          </div>
        </div>
      )}

      {toast && (
        <Toast
          message={toast.message}
          type={toast.type}
          onClose={() => setToast(null)}
        />
      )}
    </>
  );
};

export default SalesOrders;
