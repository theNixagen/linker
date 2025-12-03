import * as Minio from "minio";

export const s3Client = new Minio.Client({
  endPoint: process.env?.S3_ENDPOINT ?? "",
  accessKey: process.env.S3_ACCESS_KEY,
  secretKey: process.env.S3_SECRET_KEY,
  useSSL: false,
});

export async function createBucketIfNotExists(bucketName: string) {
  const bucketExists = await s3Client.bucketExists(bucketName);
  if (!bucketExists) {
    await s3Client.makeBucket(bucketName);
  }
}
