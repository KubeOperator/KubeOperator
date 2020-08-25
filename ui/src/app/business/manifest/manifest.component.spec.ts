import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ManifestComponent } from './manifest.component';

describe('ManifestComponent', () => {
  let component: ManifestComponent;
  let fixture: ComponentFixture<ManifestComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ManifestComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ManifestComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
