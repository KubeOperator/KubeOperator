import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ManifestAlertComponent } from './manifest-alert.component';

describe('ManifestAlertComponent', () => {
  let component: ManifestAlertComponent;
  let fixture: ComponentFixture<ManifestAlertComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ManifestAlertComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ManifestAlertComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
