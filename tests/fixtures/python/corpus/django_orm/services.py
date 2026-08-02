from django.db import models


class WarehouseStock(models.Model):
    warehouse_id = models.IntegerField()
    quantity = models.IntegerField()


class Job(models.Model):
    status = models.CharField(max_length=32)


def claim():
    job = Job.objects.filter(status="pending").first()
    job.status = "running"
    job.save()


def available():
    return WarehouseStock.objects.filter(warehouse_id=1).order_by("-quantity")


def reserve(items):
    for item in items:
        sku = WarehouseStock.objects.get(warehouse_id=item["warehouse_id"])
        _ = sku
