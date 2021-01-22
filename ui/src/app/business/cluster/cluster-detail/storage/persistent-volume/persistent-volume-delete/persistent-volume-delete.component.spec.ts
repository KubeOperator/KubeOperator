import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PersistentVolumeDeleteComponent } from './persistent-volume-delete.component';

describe('PersistentVolumeDeleteComponent', () => {
  let component: PersistentVolumeDeleteComponent;
  let fixture: ComponentFixture<PersistentVolumeDeleteComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ PersistentVolumeDeleteComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(PersistentVolumeDeleteComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
