## 📘 Disclaimer

> **This Operator is built for learning and experimentation purposes only.**
>
> The chosen API group domain (`example.com`) is used to avoid conflicts and signal non-production intent.

---

# TimeSync Operator

The **TimeSync Operator** is a Kubernetes operator designed to manage time synchronization across namespaces in a Kubernetes cluster. It automates the injection of time synchronization sidecars into pods based on namespace-level policies, ensuring consistent and accurate timekeeping across your workloads.

## Features

- **Namespace-Level Policies**: Define time synchronization policies at the namespace level using custom resources.
- **Automatic Sidecar Injection**: Automatically inject time synchronization containers into pods based on defined policies.
- **Centralized Management**: Centrally manage and enforce time synchronization across multiple namespaces.
- **Customizable Implementation**: Configure the time synchronization implementation via container image selection.
- **Secure and Observable**:
    - TLS certificates for secure communication.
    - Health and readiness probes.
    - Prometheus metrics integration.
    - Leader election support for high availability.

## Architecture

The TimeSync Operator follows the Kubernetes operator pattern and consists of two primary components:

### Controller Component

- Watches `TimeSyncPolicy` custom resources.
- Monitors namespaces that match the policy's selectors.
- Updates the `TimeSyncPolicy` status with matched namespace counts.

### Webhook Component

- Intercepts Pod creation and update requests.
- Determines if the Pod's namespace matches any `TimeSyncPolicy`.
- Injects a time synchronization sidecar container when required.

## Custom Resource Definition

The operator introduces a custom resource named `TimeSyncPolicy` that defines how time synchronization should be applied:

- **Enable or Disable**: Toggle time synchronization for specific namespaces.
- **Container Image**: Specify the container image to use for time synchronization.
- **Namespace Selection**: Define which namespaces should have time synchronization applied using label selectors.

## System Requirements

- Kubernetes v1.19 or later.
- Any Kubernetes-compliant distribution.

## License

This project is licensed under the [Apache 2 License](LICENSE).