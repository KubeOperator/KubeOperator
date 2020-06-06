import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PersistentVolumeComponent } from './persistent-volume.component';

describe('PersistentVolumeComponent', () => {
  let component: PersistentVolumeComponent;
  let fixture: ComponentFixture<PersistentVolumeComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ PersistentVolumeComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PersistentVolumeComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
