#include <stdio.h>
#include <stdlib.h>
#include <sys/mman.h>
#include <sys/resource.h>
#include <string.h>

long get_rss_kb() {
    struct rusage usage;
    getrusage(RUSAGE_SELF, &usage);
    return usage.ru_maxrss / 1024; // Convert to KB (ru_maxrss is in bytes on macOS)
}

int main() {
    size_t size = 100 * 1024 * 1024; // 100 MB
    long rss_before, rss_after_alloc, rss_after_touch, rss_after_madvise;

    printf("Testing MADV_DONTNEED on macOS\n");
    printf("==============================\n");

    // Initial RSS
    rss_before = get_rss_kb();
    printf("RSS before allocation: %ld KB\n", rss_before);

    // Allocate memory
    void *ptr = mmap(NULL, size, PROT_READ | PROT_WRITE,
                     MAP_PRIVATE | MAP_ANONYMOUS, -1, 0);
    if (ptr == MAP_FAILED) {
        perror("mmap failed");
        return 1;
    }

    rss_after_alloc = get_rss_kb();
    printf("RSS after allocation: %ld KB\n", rss_after_alloc);

    // Touch all pages to ensure they're resident
    memset(ptr, 1, size);

    rss_after_touch = get_rss_kb();
    printf("RSS after touching pages: %ld KB\n", rss_after_touch);

    // Call madvise with MADV_DONTNEED
    if (madvise(ptr, size, MADV_DONTNEED) != 0) {
        perror("madvise failed");
        munmap(ptr, size);
        return 1;
    }

    rss_after_madvise = get_rss_kb();
    printf("RSS after MADV_DONTNEED: %ld KB\n", rss_after_madvise);

    printf("\nResults:\n");
    printf("RSS reduction: %ld KB\n", rss_after_touch - rss_after_madvise);
    printf("Immediate reduction: %s\n",
           (rss_after_madvise < rss_after_touch) ? "YES" : "NO");

    // Cleanup
    munmap(ptr, size);

    return 0;
}