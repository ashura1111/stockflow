import { useEffect, useState } from 'react';
import { Trash2, Edit, Plus, TrendingDown } from 'lucide-react';
import { api, ApiError } from '../api/client';
import type { Product } from '../types/index';
import { formatCurrency, formatRelativeTime } from '../utils/format';
import Modal from '../components/modal';
import Toast from '../components/toast';
import '../App.css';

interface FormData {
  sku: string;
  name: string;
  price: number;
  stock: number;
  is_active: boolean;
}

interface StockAdjustData {
  type: 'in' | 'out';
  qty: number;
  note: string;
}

const Products = () => {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [search, setSearch] = useState('');
  const [statusFilter, setStatusFilter] = useState('');
  const [lowStockOnly, setLowStockOnly] = useState(false);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [pageSize] = useState(10);

  const [showModal, setShowModal] = useState(false);
  const [modalMode, setModalMode] = useState<'create' | 'edit'>('create');
  const [editingProduct, setEditingProduct] = useState<Product | null>(null);
  const [formData, setFormData] = useState<FormData>({
    sku: '',
    name: '',
    price: 0,
    stock: 0,
    is_active: true,
  });
  const [formError, setFormError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const [showStockModal, setShowStockModal] = useState(false);
  const [stockAdjustProduct, setStockAdjustProduct] = useState<Product | null>(null);
  const [stockAdjustData, setStockAdjustData] = useState<StockAdjustData>({
    type: 'in',
    qty: 0,
    note: '',
  });
  const [stockAdjustError, setStockAdjustError] = useState<string | null>(null);
  const [stockAdjustSubmitting, setStockAdjustSubmitting] = useState(false);

  const [showConfirmDelete, setShowConfirmDelete] = useState(false);
  const [deleteProduct, setDeleteProduct] = useState<Product | null>(null);

  const [toast, setToast] = useState<{ message: string; type: 'success' | 'error' } | null>(null);

  useEffect(() => {
    fetchProducts();
  }, [search, statusFilter, lowStockOnly, page]);

  const fetchProducts = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await api.getProducts(search, statusFilter, lowStockOnly, page, pageSize);
      setProducts(data.data);
      setTotal(data.total);
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      } else {
        setError('Failed to load products');
      }
    } finally {
      setLoading(false);
    }
  };

  const handleOpenModal = (product?: Product) => {
    if (product) {
      setModalMode('edit');
      setEditingProduct(product);
      setFormData({
        sku: product.sku,
        name: product.name,
        price: product.price,
        stock: product.stock,
        is_active: product.is_active,
      });
    } else {
      setModalMode('create');
      setEditingProduct(null);
      setFormData({
        sku: '',
        name: '',
        price: 0,
        stock: 0,
        is_active: true,
      });
    }
    setFormError(null);
    setShowModal(true);
  };

  const handleFormChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value, type } = e.target;
    if (type === 'checkbox') {
      setFormData({ ...formData, [name]: (e.target as HTMLInputElement).checked });
    } else if (type === 'number') {
      setFormData({ ...formData, [name]: parseInt(value) || 0 });
    } else {
      setFormData({ ...formData, [name]: value });
    }
  };

  const handleSubmitForm = async () => {
    setFormError(null);
    setSubmitting(true);

    try {
      if (modalMode === 'create') {
        await api.createProduct(formData);
        setToast({ message: 'Product created successfully', type: 'success' });
      } else if (editingProduct) {
        const { stock, ...updateData } = formData;
        await api.updateProduct(editingProduct.id, updateData);
        setToast({ message: 'Product updated successfully', type: 'success' });
      }
      setShowModal(false);
      await fetchProducts();
    } catch (err) {
      if (err instanceof ApiError) {
        if (err.code === 'SKU_ALREADY_EXISTS') {
          setFormError('SKU already exists');
        } else {
          setFormError(err.message);
        }
      } else {
        setFormError('Failed to save product');
      }
    } finally {
      setSubmitting(false);
    }
  };

  const handleDeleteClick = (product: Product) => {
    setDeleteProduct(product);
    setShowConfirmDelete(true);
  };

  const handleConfirmDelete = async () => {
    if (!deleteProduct) return;

    try {
      await api.deleteProduct(deleteProduct.id);
      setToast({ message: 'Product deleted successfully', type: 'success' });
      setShowConfirmDelete(false);
      await fetchProducts();
    } catch (err) {
      if (err instanceof ApiError && err.code === 'PRODUCT_HAS_HISTORY') {
        setToast({ message: 'Cannot delete product with history', type: 'error' });
      } else {
        setToast({ message: 'Failed to delete product', type: 'error' });
      }
      setShowConfirmDelete(false);
    }
  };

  const handleOpenStockModal = (product: Product) => {
    setStockAdjustProduct(product);
    setStockAdjustData({ type: 'in', qty: 0, note: '' });
    setStockAdjustError(null);
    setShowStockModal(true);
  };

  const handleStockAdjustChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value, type } = e.target;
    if (type === 'number') {
      setStockAdjustData({ ...stockAdjustData, [name]: parseInt(value) || 0 });
    } else {
      setStockAdjustData({ ...stockAdjustData, [name]: value });
    }
  };

  const handleSubmitStockAdjust = async () => {
    if (!stockAdjustProduct) return;
    setStockAdjustError(null);
    setStockAdjustSubmitting(true);

    try {
      await api.adjustStock(stockAdjustProduct.id, stockAdjustData);
      setToast({ message: 'Stock adjusted successfully', type: 'success' });
      setShowStockModal(false);
      await fetchProducts();
    } catch (err) {
      if (err instanceof ApiError) {
        if (err.code === 'INSUFFICIENT_STOCK') {
          setStockAdjustError(`Insufficient stock. Available: ${err.details.available}`);
        } else {
          setStockAdjustError(err.message);
        }
      } else {
        setStockAdjustError('Failed to adjust stock');
      }
    } finally {
      setStockAdjustSubmitting(false);
    }
  };

  const totalPages = Math.ceil(total / pageSize);

  return (
    <>
      <div className="page-header">
        <h2>Products</h2>
        <p>Manage your product inventory</p>
      </div>

      <div className="toolbar" style={{ marginBottom: '20px' }}>
        <input
          type="text"
          className="search-input"
          placeholder="Search by SKU or name..."
          value={search}
          onChange={(e) => {
            setSearch(e.target.value);
            setPage(1);
          }}
        />
        <select
          className="form-control"
          value={statusFilter}
          onChange={(e) => {
            setStatusFilter(e.target.value);
            setPage(1);
          }}
          style={{ flex: '0 1 120px' }}
        >
          <option value="">All</option>
          <option value="active">Active</option>
          <option value="inactive">Inactive</option>
        </select>
        <label style={{ display: 'flex', alignItems: 'center', gap: '6px', cursor: 'pointer' }}>
          <input
            type="checkbox"
            checked={lowStockOnly}
            onChange={(e) => {
              setLowStockOnly(e.target.checked);
              setPage(1);
            }}
          />
          Low Stock
        </label>
        <button className="btn btn-primary" onClick={() => handleOpenModal()}>
          <Plus size={16} /> Add Product
        </button>
      </div>

      {error && <div className="error-message">{error}</div>}

      {loading ? (
        <div className="loading">Loading products...</div>
      ) : products.length === 0 ? (
        <div className="empty-state">
          <p>No products found</p>
          <small>Create your first product to get started</small>
        </div>
      ) : (
        <div className="card">
          <div className="table-container">
            <table>
              <thead>
                <tr>
                  <th>SKU</th>
                  <th>Name</th>
                  <th>Price</th>
                  <th>Stock</th>
                  <th>Status</th>
                  <th>Last Movement</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {products.map((product) => (
                  <tr key={product.id}>
                    <td>{product.sku}</td>
                    <td>{product.name}</td>
                    <td>{formatCurrency(product.price)}</td>
                    <td>
                      <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                        {product.stock}
                        {product.stock <= 10 && <TrendingDown size={16} color="#d97706" />}
                      </div>
                    </td>
                    <td>
                      <span className={`badge ${product.is_active ? 'badge-success' : 'badge-danger'}`}>
                        {product.is_active ? 'Active' : 'Inactive'}
                      </span>
                    </td>
                    <td style={{ fontSize: '13px', color: '#64748b' }}>
                      {product.last_movement_at ? formatRelativeTime(product.last_movement_at) : 'Never'}
                    </td>
                    <td>
                      <div style={{ display: 'flex', gap: '8px' }}>
                        <button className="btn btn-secondary" onClick={() => handleOpenStockModal(product)}>
                          Stock
                        </button>
                        <button className="btn btn-secondary" onClick={() => handleOpenModal(product)}>
                          <Edit size={14} />
                        </button>
                        <button className="btn btn-danger" onClick={() => handleDeleteClick(product)}>
                          <Trash2 size={14} />
                        </button>
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
        isOpen={showModal}
        onClose={() => setShowModal(false)}
        title={modalMode === 'create' ? 'Create Product' : 'Edit Product'}
        onSubmit={handleSubmitForm}
        submitLabel={modalMode === 'create' ? 'Create' : 'Update'}
        isLoading={submitting}
      >
        {formError && <div className="error-message" style={{ marginBottom: '12px' }}>{formError}</div>}

        <div className="form-group">
          <label>SKU</label>
          <input
            type="text"
            className="form-control"
            name="sku"
            value={formData.sku}
            onChange={handleFormChange}
            disabled={modalMode === 'edit'}
            placeholder="e.g., CRM-001"
          />
        </div>

        <div className="form-group">
          <label>Product Name</label>
          <input
            type="text"
            className="form-control"
            name="name"
            value={formData.name}
            onChange={handleFormChange}
            placeholder="e.g., Keramik Lantai 40x40 cm"
          />
        </div>

        <div className="form-group">
          <label>Price (IDR)</label>
          <input
            type="number"
            className="form-control"
            name="price"
            value={formData.price}
            onChange={handleFormChange}
            placeholder="0"
            min="0"
          />
        </div>

        {modalMode === 'create' && (
          <div className="form-group">
            <label>Initial Stock</label>
            <input
              type="number"
              className="form-control"
              name="stock"
              value={formData.stock}
              onChange={handleFormChange}
              placeholder="0"
              min="0"
            />
          </div>
        )}

        <div className="form-group">
          <label style={{ display: 'flex', gap: '6px', alignItems: 'center', cursor: 'pointer' }}>
            <input
              type="checkbox"
              name="is_active"
              checked={formData.is_active}
              onChange={handleFormChange}
            />
            Active
          </label>
        </div>
      </Modal>

      <Modal
        isOpen={showStockModal}
        onClose={() => setShowStockModal(false)}
        title={`Adjust Stock - ${stockAdjustProduct?.sku}`}
        onSubmit={handleSubmitStockAdjust}
        submitLabel="Adjust"
        isLoading={stockAdjustSubmitting}
      >
        {stockAdjustError && <div className="error-message" style={{ marginBottom: '12px' }}>{stockAdjustError}</div>}

        <div className="form-group">
          <label>Current Stock: {stockAdjustProduct?.stock}</label>
        </div>

        <div className="form-group">
          <label>Type</label>
          <select
            className="form-control"
            name="type"
            value={stockAdjustData.type}
            onChange={handleStockAdjustChange}
          >
            <option value="in">Stock In</option>
            <option value="out">Stock Out</option>
          </select>
        </div>

        <div className="form-group">
          <label>Quantity</label>
          <input
            type="number"
            className="form-control"
            name="qty"
            value={stockAdjustData.qty}
            onChange={handleStockAdjustChange}
            placeholder="0"
            min="1"
          />
        </div>

        <div className="form-group">
          <label>Note</label>
          <input
            type="text"
            className="form-control"
            name="note"
            value={stockAdjustData.note}
            onChange={handleStockAdjustChange}
            placeholder="Optional note"
          />
        </div>
      </Modal>

      {showConfirmDelete && deleteProduct && (
        <div className="modal-overlay" onClick={() => setShowConfirmDelete(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h3>Confirm Delete</h3>
            </div>
            <div className="modal-body">
              <p>Are you sure you want to delete product <strong>{deleteProduct.sku}</strong>?</p>
              <p style={{ fontSize: '13px', color: '#64748b', marginTop: '12px' }}>
                This action cannot be undone.
              </p>
            </div>
            <div className="modal-footer">
              <button className="btn btn-secondary" onClick={() => setShowConfirmDelete(false)}>
                Cancel
              </button>
              <button className="btn btn-danger" onClick={handleConfirmDelete}>
                Delete
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

export default Products;
