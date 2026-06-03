import { useEffect, useState } from "react";
import { Link, useParams, useNavigate } from "react-router-dom";
import {
  getProduct, deleteProduct, listSubstitutes, addSubstitute, removeSubstitute,
  listProducts, listCategories,
  type Product, type GenericSubstitution, type ProductCategory,
} from "../api/products";
import { useToast } from "../context/ToastContext";

export default function ProductDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [product, setProduct] = useState<Product | null>(null);
  const [categories, setCategories] = useState<ProductCategory[]>([]);
  const [substitutes, setSubstitutes] = useState<GenericSubstitution[]>([]);
  const [allProducts, setAllProducts] = useState<Product[]>([]);
  const [newSub, setNewSub] = useState("");
  const [newSubNotes, setNewSubNotes] = useState("");
  const [loading, setLoading] = useState(true);
  const { showToast } = useToast();

  const load = () => {
    if (!id) return;
    setLoading(true);
    Promise.all([
      getProduct(id),
      listCategories(),
      listSubstitutes(id),
      listProducts(1, 200),
    ]).then(([p, cats, subs, all]) => {
      setProduct(p);
      setCategories(cats);
      setSubstitutes(subs);
      setAllProducts(all.data);
    }).catch(console.error).finally(() => setLoading(false));
  };

  useEffect(load, [id]);

  const catMap = Object.fromEntries(categories.map(c => [c.id, c.name]));
  const prodMap = Object.fromEntries(allProducts.map(p => [p.id, p.name]));

  const handleDelete = async () => {
    if (!confirm("Delete this product?")) return;
    try {
      await deleteProduct(id!);
      navigate("/products");
    } catch (err) {
      showToast(err instanceof Error ? err.message : "Delete failed", "error");
    }
  };

  const handleAddSub = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newSub) return;
    try {
      await addSubstitute(id!, newSub, newSubNotes || undefined);
      setNewSub("");
      setNewSubNotes("");
      setSubstitutes(await listSubstitutes(id!));
    } catch (err) {
      showToast(err instanceof Error ? err.message : "Failed", "error");
    }
  };

  const handleRemoveSub = async (subId: string) => {
    try {
      await removeSubstitute(id!, subId);
      setSubstitutes(substitutes.filter(s => s.id !== subId));
    } catch (err) {
      showToast(err instanceof Error ? err.message : "Failed", "error");
    }
  };

  if (loading) return <p>Loading...</p>;
  if (!product) return <p>Product not found</p>;

  return (
    <div>
      <div className="page-header">
        <h1>{product.name}</h1>
        <div>
          <Link to={`/products/${id}/edit`} className="btn">Edit</Link>
          <button onClick={handleDelete} style={{ marginLeft: 8 }}>Delete</button>
        </div>
      </div>

      <table>
        <tbody>
          <tr><td><strong>Brand</strong></td><td>{product.brand_name || "—"}</td></tr>
          <tr><td><strong>Generic</strong></td><td>{product.generic_name || "—"}</td></tr>
          <tr><td><strong>Manufacturer</strong></td><td>{product.manufacturer || "—"}</td></tr>
          <tr><td><strong>Category</strong></td><td>{catMap[product.category_id] || "—"}</td></tr>
          <tr><td><strong>Classification</strong></td><td>{product.classification || "—"}</td></tr>
          <tr><td><strong>Strength</strong></td><td>{product.strength || "—"}</td></tr>
          <tr><td><strong>Form</strong></td><td>{product.form || "—"}</td></tr>
          <tr><td><strong>Unit</strong></td><td>{product.unit_of_measure || "—"}</td></tr>
          <tr><td><strong>Barcode</strong></td><td>{product.barcode || "—"}</td></tr>
          <tr><td><strong>NDC</strong></td><td>{product.ndc || "—"}</td></tr>
          <tr><td><strong>Reorder Level</strong></td><td>{product.reorder_level}</td></tr>
          <tr><td><strong>Status</strong></td><td>{product.is_active ? "Active" : "Inactive"}</td></tr>
          <tr><td><strong>Description</strong></td><td>{product.description || "—"}</td></tr>
        </tbody>
      </table>

      <h2 style={{ marginTop: 32 }}>Generic Substitutes</h2>
      <form onSubmit={handleAddSub} style={{ display: "flex", gap: 8, maxWidth: "100%", marginBottom: 12 }}>
        <select value={newSub} onChange={e => setNewSub(e.target.value)} style={{ flex: 1 }}>
          <option value="">Select product...</option>
          {allProducts
            .filter(p => p.id !== id)
            .map(p => <option key={p.id} value={p.id}>{p.name} ({p.generic_name || p.brand_name || "—"})</option>)
          }
        </select>
        <input value={newSubNotes} onChange={e => setNewSubNotes(e.target.value)} placeholder="Notes" />
        <button type="submit">Add</button>
      </form>
      {substitutes.length === 0 ? <p>No substitutes defined.</p> : (
        <table>
          <thead><tr><th>Product</th><th>Notes</th><th>Actions</th></tr></thead>
          <tbody>
            {substitutes.map(s => (
              <tr key={s.id}>
                <td>{prodMap[s.substitute_product_id] || s.substitute_product_id}</td>
                <td>{s.notes || "—"}</td>
                <td><button onClick={() => handleRemoveSub(s.id)}>Remove</button></td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}
