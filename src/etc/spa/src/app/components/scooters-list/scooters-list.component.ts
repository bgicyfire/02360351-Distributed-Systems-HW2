import {Component, OnInit} from '@angular/core';
import {ScootersService} from '../../services/scooters.service';
import {map, Observable} from 'rxjs';
import {Scooter} from '../../dtos/scooter';
import {AsyncPipe, NgClass, NgForOf, NgIf} from '@angular/common';
import {ToastService} from '../../services/toast.service';

@Component({
  selector: 'app-scooters-list',
  imports: [
    AsyncPipe,
    NgForOf,
    NgClass,
    NgIf,
  ],
  templateUrl: './scooters-list.component.html',
  styleUrl: './scooters-list.component.sass'
})
export class ScootersListComponent implements OnInit {
  protected scooters!: Observable<Scooter[]>;

  constructor(
    private scootersService: ScootersService,
    private toastService: ToastService) {
  }

  public ngOnInit(): void {
    this.refresh();
  }

  protected refresh(): void {
    this.scooters = this.scootersService.getScootersList();
  }

  protected createScooter(id: string): void {
    this.scootersService.createScooter(id)
      .subscribe(a => {
        this.refresh();
      });
  }

  protected reserveScooter(id: string): void {
    this.scootersService.reserveScooter(id)
      .subscribe(a => {
        this.refresh();
      });
  }

  protected releaseScooter(id: string, reservationId: string, rideDistance: number): void {
    this.scootersService.releaseScooter(id, reservationId, rideDistance)
      .subscribe(a => {
        this.refresh();
      });
  }

  protected readonly parseInt = parseInt;
}
