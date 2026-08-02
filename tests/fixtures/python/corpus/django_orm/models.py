from django.db import models


class WarehouseStock(models.Model):
    warehouse_id = models.IntegerField()
    quantity = models.IntegerField()


class Job(models.Model):
    status = models.CharField(max_length=32)
