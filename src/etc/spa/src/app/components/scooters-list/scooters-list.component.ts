import {Component, OnInit} from '@angular/core';
import {ScootersService} from '../../services/scooters.service';
import {Observable} from 'rxjs';
import {Scooter} from '../../dtos/scooter';
import {AsyncPipe, NgClass, NgForOf, NgIf} from '@angular/common';

@Component({
  selector: 'app-scooters-list',
  imports: [
    AsyncPipe,
    NgForOf,
    NgClass,
    NgIf
  ],
  templateUrl: './scooters-list.component.html',
  styleUrl: './scooters-list.component.sass'
})
export class ScootersListComponent implements OnInit {
  protected scooters!: Observable<Scooter[]>;

  constructor(private scootersService: ScootersService) {
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
        alert('Scooter was created successfully');
        this.refresh();
      });
  }

  protected reserveScooter(id: string): void {
    this.scootersService.reserveScooter(id)
      .subscribe(a => {
        alert('Scooter was reserved successfully');
        this.refresh();
      });
  }

  protected releaseScooter(id: string, reservationId: string, rideDistance: number): void {
    this.scootersService.releaseScooter(id, reservationId, rideDistance)
      .subscribe(a => {
        alert('Scooter was released successfully');
        this.refresh();
      });
  }

  protected readonly parseInt = parseInt;
}
