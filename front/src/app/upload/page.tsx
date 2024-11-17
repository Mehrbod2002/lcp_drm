"use client";

import { useForm } from "react-hook-form";
import axios from "axios";
import { signOut } from "next-auth/react";
import { useRouter } from "next/navigation";
import { useState } from "react";

type UploadForm = {
  username: string;
  password: string;
  copy: number;
  print: number;
  end: number;
  file: FileList;
};

export default function UploadPage() {
  const { register, handleSubmit, reset } = useForm<UploadForm>();
  const router = useRouter();
  const [downloadLink, setDownloadLink] = useState<string | null>(null);

  const onSubmit = async (data: UploadForm) => {
    const formData = new FormData();
    formData.append("username", data.username);
    formData.append("password", data.password);
    formData.append("copy", data.copy.toString());
    formData.append("print", data.print.toString());
    formData.append("end", data.end.toString());
    formData.append("file", data.file[0]);

    try {
      const res = await axios.post(
        "https://bookadd.ir/backend/upload",
        formData,
        { headers: { "Content-Type": "multipart/form-data" } }
      );

      alert("File uploaded successfully!");
      setDownloadLink(`https://bookadd.ir${res.data.download_link}`);
      reset();
    } catch (err) {
      console.error(err);
      alert("Upload failed!");
    }
  };

  return (
    <div className="min-h-screen flex flex-col bg-gray-50">
      {/* Header */}
      <header className="bg-blue-600 text-white py-4 shadow-md">
        <div className="container mx-auto flex justify-between items-center px-4">
          <h1 className="text-2xl font-bold">Upload Center</h1>
          <button
            onClick={() => {
              signOut();
              router.push("/login");
            }}
            className="bg-red-500 hover:bg-red-700 text-white font-medium py-2 px-4 rounded"
          >
            Logout
          </button>
        </div>
      </header>

      {/* Main Content */}
      <main className="flex-grow flex flex-col items-center justify-center px-4">
        <form
          onSubmit={handleSubmit(onSubmit)}
          className="bg-white shadow-lg rounded-lg p-8 w-full max-w-md space-y-6"
        >
          <h2 className="text-2xl font-bold text-center text-gray-800">
            Upload Your File
          </h2>

          {/* Username */}
          <div>
            <label className="block text-gray-700 font-medium mb-1">Username</label>
            <input
              {...register("username")}
              placeholder="Enter your username"
              className="shadow-sm appearance-none border rounded w-full py-2 px-3 text-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          {/* Password */}
          <div>
            <label className="block text-gray-700 font-medium mb-1">Password</label>
            <input
              type="password"
              {...register("password")}
              placeholder="Enter your password"
              className="shadow-sm appearance-none border rounded w-full py-2 px-3 text-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          {/* Copy Limit */}
          <div>
            <label className="block text-gray-700 font-medium mb-1">Copy Limit</label>
            <input
              type="number"
              {...register("copy", { valueAsNumber: true })}
              placeholder="Enter copy limit"
              className="shadow-sm appearance-none border rounded w-full py-2 px-3 text-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          {/* Print Limit */}
          <div>
            <label className="block text-gray-700 font-medium mb-1">Print Limit</label>
            <input
              type="number"
              {...register("print", { valueAsNumber: true })}
              placeholder="Enter print limit"
              className="shadow-sm appearance-none border rounded w-full py-2 px-3 text-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          {/* End Limit */}
          <div>
            <label className="block text-gray-700 font-medium mb-1">End Limit</label>
            <input
              type="number"
              {...register("end", { valueAsNumber: true })}
              placeholder="Enter end limit"
              className="shadow-sm appearance-none border rounded w-full py-2 px-3 text-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          {/* File Upload */}
          <div>
            <label className="block text-gray-700 font-medium mb-1">Upload File</label>
            <input
              type="file"
              {...register("file")}
              className="block w-full text-gray-700 border border-gray-300 rounded-lg cursor-pointer bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          {/* Submit Button */}
          <button
            type="submit"
            className="w-full bg-blue-500 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
          >
            Submit
          </button>
        </form>

        {/* Display Download Link */}
        {downloadLink && (
          <div className="mt-8 text-center">
            <p className="text-lg font-medium text-gray-800">
              File encrypted successfully! Download it here:
            </p>
            <a
              href={downloadLink}
              download
              className="text-blue-600 underline hover:text-blue-800"
            >
              {downloadLink}
            </a>
          </div>
        )}
      </main>
    </div>
  );
}
