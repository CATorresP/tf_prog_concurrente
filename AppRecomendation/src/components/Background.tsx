import React, { useEffect, useRef } from "react";

export const DynamicBackground: React.FC = () => {
  const canvasRef = useRef<HTMLCanvasElement | null>(null);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (canvas) {
      const context = canvas.getContext("2d");
      if (!context) return;

      const particles: {
        x: number;
        y: number;
        size: number;
        dx: number;
        dy: number;
      }[] = [];
      const numParticles = 100;
      const particleSize = 2;

      const resizeCanvas = () => {
        canvas.width = window.innerWidth;
        canvas.height = window.innerHeight;
      };

      // Initialize particles
      for (let i = 0; i < numParticles; i++) {
        particles.push({
          x: Math.random() * canvas.width,
          y: Math.random() * canvas.height,
          size: particleSize,
          dx: (Math.random() - 0.5) * 2,
          dy: (Math.random() - 0.5) * 2,
        });
      }

      const animate = () => {
        if (!canvas || !context) return;
        context.clearRect(0, 0, canvas.width, canvas.height);

        particles.forEach((particle) => {
          context.beginPath();
          context.arc(particle.x, particle.y, particle.size, 0, Math.PI * 2);
          context.fillStyle = "rgba(255, 255, 255, 0.7)";
          context.fill();
          context.closePath();

          // Move particle
          particle.x += particle.dx;
          particle.y += particle.dy;

          // Check boundaries
          if (
            particle.x - particle.size < 0 ||
            particle.x + particle.size > canvas.width
          ) {
            particle.dx = -particle.dx;
          }
          if (
            particle.y - particle.size < 0 ||
            particle.y + particle.size > canvas.height
          ) {
            particle.dy = -particle.dy;
          }
        });

        requestAnimationFrame(animate);
      };

      resizeCanvas();
      animate();

      window.addEventListener("resize", resizeCanvas);

      return () => {
        window.removeEventListener("resize", resizeCanvas);
      };
    }
  }, []);

  return (
    <canvas
      ref={canvasRef}
      style={{
        display: "block",
        position: "fixed",
        backgroundColor: "black",
        top: 0,
        left: 0,
        zIndex: -1,
      }}
    />
  );
};
