import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import axios from "axios";
import { Button, Card } from "flowbite-react";
import { toast } from "react-toastify";

interface Vendor {
  vendor_id: string;
  nama_vendor: string;
  alamat: string;
  status: string;
}

export default function VendorDetail() {
  const [vendor, setVendor] = useState<Vendor | null>(null);
  const { id } = useParams();
  const navigate = useNavigate();
  const token = localStorage.getItem("token");

  useEffect(() => {
    const fetchDetail = async () => {
      try {
        const res = await axios.get(`http://localhost:6969/vendor/${id}`, {
          headers: { Authorization: `Bearer ${token}` },
        });
        setVendor(res.data.data);
      } catch (error) {
        toast.error("Gagal mengambil data vendor");
      }
    };
    fetchDetail();
  }, [id, token]);

  const handleDelete = async () => {
    if (!confirm("Yakin ingin menghapus vendor ini?")) return;
    try {
      await axios.delete(`http://localhost:6969/vendor-kontak/delete/${id}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      toast.success("Vendor berhasil dihapus");
      navigate("/vendor");
    } catch {
      toast.error("Gagal menghapus vendor");
    }
  };

  if (!vendor) return <div>Loading...</div>;

  return (
    <div className="p-6 max-w-xl mx-auto">
      <Card>
        <h1 className="text-2xl font-bold mb-4">Detail Vendor</h1>
        <div className="space-y-2">
          <p><strong>Nama:</strong> {vendor.nama_vendor}</p>
          <p><strong>Alamat:</strong> {vendor.alamat}</p>
          <p><strong>Status:</strong> {vendor.status}</p>
        </div>
        <div className="flex gap-4 mt-6">
          <Button color="blue" onClick={() => navigate(`/vendor/edit/${vendor.vendor_id}`)}>Edit</Button>
          <Button color="failure" onClick={handleDelete}>Delete</Button>
        </div>
      </Card>
    </div>
  );
}
